package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	f "github.com/gruntwork-io/terratest/modules/files"
	ts "github.com/gruntwork-io/terratest/modules/test-structure"
	cp "github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

// Global test variables
var (
	subscriptionID     = "<private_value>"
	location           = "CanadaCentral"
	clientName         = "client"
	environment        = "test"
	stack              = "stack"
	suffixLength       = 8
	setupTerraformDir  = "terraform"
	moduleTerraformDir = "../"
	testModuleDir        = "module/"
	testSetupDir         = "terraform/"
	testModuleTerraformOptionsDir = "subnetTerraformOptions/"
	testSetupTerraformOptionsDir = "setupTerraformOptions/"
)

type SubnetTestData struct {
	vNetRgName string
	// Variables for vnet
	vNetCidr string
	vNetName string
	// Variables for subnet
	subnetCidr string
	privateEndpointEnabled bool
	privateLinkServiceEnabled bool
	expectedSubnetName string
}

type SubnetWithServiceEndpointsTestData struct {
	vNetRgName string
	// Variables for vnet
	vNetCidr string
	vNetName string
	// Variables for subnet
	subnetCidr string
	privateEndpointEnabled bool
	privateLinkServiceEnabled bool
	expectedSubnetName string
	serviceEndpoints []string
}

func TestSubnetWithDefaultConfigs(t *testing.T) {
	testRootDir := "TestSubnetWithDefaultConfigs/"

	// Uncomment any of the following lines to skip that test stage
	// os.Setenv("SKIP_setup_" + testRootDir, "true")
	// os.Setenv("SKIP_deploy_" + testRootDir, "true")
	// os.Setenv("SKIP_validate_" + testRootDir, "true")
	// os.Setenv("SKIP_teardown_" + testRootDir, "true")

	t.Parallel() // Remove to test serially

	nameSuffix := GetNameSuffix(t, testRootDir)

	testData := SubnetTestData {
		vNetRgName: fmt.Sprintf("rg-snet-unit-test-%s", nameSuffix),
		vNetCidr: "10.0.0.0/16",
		vNetName: fmt.Sprintf("vnet-snet-unit-test-%s", nameSuffix),
		subnetCidr: "10.0.0.0/24",
		privateEndpointEnabled: false,
		privateLinkServiceEnabled: false,
		expectedSubnetName: fmt.Sprintf("snet-stack-client-test-%s", nameSuffix),
	}

	SubnetWithDefaultConfigs(t, testRootDir, nameSuffix, testData)
}


func SubnetWithDefaultConfigs(t *testing.T, testRootDir string, nameSuffix string, testData SubnetTestData) {
	// At the end of the test, clean up resources.
	defer ts.RunTestStage(t, "teardown_" + testRootDir, func() {
		TearDown(t, testRootDir)
	})

	ts.RunTestStage(t, "setup_" + testRootDir, func() {
		// If state files exist, clean up resources
		TearDown(t, testRootDir)
		CopyTerraformFolder(setupTerraformDir, fmt.Sprintf("%s%s", testRootDir, testSetupDir))
		ts.SaveString(t, testRootDir, "nameSuffix", nameSuffix)
		
		setupTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: fmt.Sprintf("%s%s", testRootDir, testSetupDir),
			Vars: map[string]interface{}{
				"config": map[string]interface{}{
					"location":            location,
					"resource_group_name": testData.vNetRgName,
					"vnet_name":           testData.vNetName,
					"address_space":       []string{testData.vNetCidr},
				},
			},
		})

		ts.SaveTerraformOptions(t, fmt.Sprintf("%s%s", testRootDir, testSetupTerraformOptionsDir), setupTerraformOptions)
		terraform.InitAndApply(t, setupTerraformOptions)
	})

	ts.RunTestStage(t, "deploy_" + testRootDir, func() {
		CopyTerraformFolder(moduleTerraformDir, fmt.Sprintf("%s%s", testRootDir, testModuleDir))

		subnetTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: fmt.Sprintf("%s%s", testRootDir, testModuleDir),
			Vars: map[string]interface{}{
				"vnet_resource_group_name":     testData.vNetRgName,
				"vnet_name":                    testData.vNetName,
				"name_suffix":                  nameSuffix,
				"stack":                        stack,
				"environment":                  environment,
				"client_name":                  clientName,
				"subnet_cidr_list":             []string{testData.subnetCidr},
				"private_endpoint_enabled":     testData.privateEndpointEnabled,
				"private_link_service_enabled": testData.privateLinkServiceEnabled,
			},
		})

		ts.SaveTerraformOptions(t, fmt.Sprintf("%s%s", testRootDir, testModuleTerraformOptionsDir), subnetTerraformOptions)

		terraform.InitAndApply(t, subnetTerraformOptions)
	})

	ts.RunTestStage(t, "validate_" + testRootDir, func() {		
		vNetSubnets := azure.GetVirtualNetworkSubnets(t, testData.vNetName, testData.vNetRgName, subscriptionID)
		// Ensure subnet is present in the virtual network
		assert.NotNil(t, vNetSubnets[testData.expectedSubnetName])
		// Ensure subnet has the correct address space
		assert.Equal(t, testData.subnetCidr, vNetSubnets[testData.expectedSubnetName])
		// Get the subnet and store in object
		deployedSubnet, err := azure.GetSubnetE(testData.expectedSubnetName, testData.vNetName, testData.vNetRgName, subscriptionID)
		assert.Nil(t, err)
		// Get the subnet's properties
		deployedSubnetProperties := deployedSubnet.SubnetPropertiesFormat
		// Ensure that private endpoint link, and endpoint policies are disabled
		assert.Equal(t, *deployedSubnetProperties.PrivateEndpointNetworkPolicies, "Disabled")
		assert.Equal(t, *deployedSubnetProperties.PrivateLinkServiceNetworkPolicies, "Disabled")
	})
}

func TestSubnetWithPrivatePolicies(t *testing.T) {
	testRootDir := "TestSubnetWithPrivatePolicies/"

	// Uncomment any of the following lines to skip that test stage
	// os.Setenv("SKIP_setup_" + testRootDir, "true")
	// os.Setenv("SKIP_deploy_" + testRootDir, "true")
	// os.Setenv("SKIP_validate_" + testRootDir, "true")
	// os.Setenv("SKIP_teardown_" + testRootDir, "true")

	t.Parallel() // Remove to test serially

	nameSuffix := GetNameSuffix(t, testRootDir)

	testData := SubnetTestData {
		vNetRgName: fmt.Sprintf("rg-snet-unit-test-%s", nameSuffix),
		vNetCidr: "10.2.0.0/16",
		vNetName: fmt.Sprintf("vnet-snet-unit-test-%s", nameSuffix),
		subnetCidr: "10.2.0.0/24",
		privateEndpointEnabled: true,
		privateLinkServiceEnabled: true,
		expectedSubnetName: fmt.Sprintf("snet-stack-client-test-%s", nameSuffix),
	}

	SubnetWithPrivatePolicies(t, testRootDir, nameSuffix, testData)
}

func SubnetWithPrivatePolicies(t *testing.T, testRootDir string, nameSuffix string, testData SubnetTestData) {
	// At the end of the test, clean up resources.
	defer ts.RunTestStage(t, "teardown_" + testRootDir, func() {
		TearDown(t, testRootDir)
	})

	ts.RunTestStage(t, "setup_" + testRootDir, func() {
		// If state files exist, clean up resources
		TearDown(t, testRootDir)
		CopyTerraformFolder(setupTerraformDir, fmt.Sprintf("%s%s", testRootDir, testSetupDir))
		ts.SaveString(t, testRootDir, "nameSuffix", nameSuffix)
		
		setupTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: fmt.Sprintf("%s%s", testRootDir, testSetupDir),
			Vars: map[string]interface{}{
				"config": map[string]interface{}{
					"location":            location,
					"resource_group_name": testData.vNetRgName,
					"vnet_name":           testData.vNetName,
					"address_space":       []string{testData.vNetCidr},
				},
			},
		})

		ts.SaveTerraformOptions(t, fmt.Sprintf("%s%s", testRootDir, testSetupTerraformOptionsDir), setupTerraformOptions)
		terraform.InitAndApply(t, setupTerraformOptions)
	})

	ts.RunTestStage(t, "deploy_" + testRootDir, func() {
		CopyTerraformFolder(moduleTerraformDir, fmt.Sprintf("%s%s", testRootDir, testModuleDir))

		subnetTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: fmt.Sprintf("%s%s", testRootDir, testModuleDir),
			Vars: map[string]interface{}{
				"vnet_resource_group_name":     testData.vNetRgName,
				"vnet_name":                    testData.vNetName,
				"name_suffix":                  nameSuffix,
				"stack":                        stack,
				"environment":                  environment,
				"client_name":                  clientName,
				"subnet_cidr_list":             []string{testData.subnetCidr},
				"private_endpoint_enabled":     testData.privateEndpointEnabled,
				"private_link_service_enabled": testData.privateLinkServiceEnabled,
			},
		})

		ts.SaveTerraformOptions(t, fmt.Sprintf("%s%s", testRootDir, testModuleTerraformOptionsDir), subnetTerraformOptions)

		terraform.InitAndApply(t, subnetTerraformOptions)
	})

	ts.RunTestStage(t, "validate_" + testRootDir, func() {
		vNetSubnets := azure.GetVirtualNetworkSubnets(t, testData.vNetName, testData.vNetRgName, subscriptionID)
		// Ensure subnet is present in the virtual network
		assert.NotNil(t, vNetSubnets[testData.expectedSubnetName])
		// Ensure subnet has the correct address space
		assert.Equal(t, testData.subnetCidr, vNetSubnets[testData.expectedSubnetName])
		// Get the subnet and store in object
		deployedSubnet, err := azure.GetSubnetE(testData.expectedSubnetName, testData.vNetName, testData.vNetRgName, subscriptionID)
		assert.Nil(t, err)
		// Get the subnet's properties
		deployedSubnetProperties := deployedSubnet.SubnetPropertiesFormat
		// Ensure that private endpoint link, and endpoint policies are enabled
		assert.Equal(t, *deployedSubnetProperties.PrivateEndpointNetworkPolicies, "Enabled")
		assert.Equal(t, *deployedSubnetProperties.PrivateLinkServiceNetworkPolicies, "Enabled")
	})
}

func TestSubnetWithServiceEndpoints(t *testing.T) {
	testRootDir := "TestSubnetWithServiceEndpoints/"

	// Uncomment any of the following lines to skip that test stage
	// os.Setenv("SKIP_setup_" + testRootDir, "true")
	// os.Setenv("SKIP_deploy_" + testRootDir, "true")
	// os.Setenv("SKIP_validate_" + testRootDir, "true")
	// os.Setenv("SKIP_teardown_" + testRootDir, "true")

	t.Parallel() // Remove to test serially

	nameSuffix := GetNameSuffix(t, testRootDir)

	testData := SubnetWithServiceEndpointsTestData {
		vNetRgName: fmt.Sprintf("rg-snet-unit-test-%s", nameSuffix),
		vNetCidr: "10.4.0.0/16",
		vNetName: fmt.Sprintf("vnet-snet-unit-test-%s", nameSuffix),
		subnetCidr: "10.4.0.0/24",
		privateEndpointEnabled: false,
		privateLinkServiceEnabled: false,
		expectedSubnetName: fmt.Sprintf("snet-stack-client-test-%s", nameSuffix),
		serviceEndpoints: []string{"Microsoft.Storage", "Microsoft.Sql", "Microsoft.ServiceBus"},
	}

	SubnetWithServiceEndpoints(t, testRootDir, nameSuffix, testData)
}

func SubnetWithServiceEndpoints(t *testing.T, testRootDir string, nameSuffix string, testData SubnetWithServiceEndpointsTestData) {
	// At the end of the test, clean up resources.
	defer ts.RunTestStage(t, "teardown_" + testRootDir, func() {
		TearDown(t, testRootDir)
	})

	ts.RunTestStage(t, "setup_" + testRootDir, func() {
		// If state files exist, clean up resources
		TearDown(t, testRootDir)
		CopyTerraformFolder(setupTerraformDir, fmt.Sprintf("%s%s", testRootDir, testSetupDir))
		ts.SaveString(t, testRootDir, "nameSuffix", nameSuffix)
		
		setupTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: fmt.Sprintf("%s%s", testRootDir, testSetupDir),
			Vars: map[string]interface{}{
				"config": map[string]interface{}{
					"location":            location,
					"resource_group_name": testData.vNetRgName,
					"vnet_name":           testData.vNetName,
					"address_space":       []string{testData.vNetCidr},
				},
			},
		})

		ts.SaveTerraformOptions(t, fmt.Sprintf("%s%s", testRootDir, testSetupTerraformOptionsDir), setupTerraformOptions)
		terraform.InitAndApply(t, setupTerraformOptions)
	})

	ts.RunTestStage(t, "deploy_" + testRootDir, func() {
		CopyTerraformFolder(moduleTerraformDir, fmt.Sprintf("%s%s", testRootDir, testModuleDir))

		subnetTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: fmt.Sprintf("%s%s", testRootDir, testModuleDir),
			Vars: map[string]interface{}{
				"vnet_resource_group_name":     testData.vNetRgName,
				"vnet_name":                    testData.vNetName,
				"name_suffix":                  nameSuffix,
				"stack":                        stack,
				"environment":                  environment,
				"client_name":                  clientName,
				"subnet_cidr_list":             []string{testData.subnetCidr},
				"private_endpoint_enabled":     testData.privateEndpointEnabled,
				"private_link_service_enabled": testData.privateLinkServiceEnabled,
				"service_endpoints":            testData.serviceEndpoints,
			},
		})

		ts.SaveTerraformOptions(t, fmt.Sprintf("%s%s", testRootDir, testModuleTerraformOptionsDir), subnetTerraformOptions)

		terraform.InitAndApply(t, subnetTerraformOptions)
	})

	ts.RunTestStage(t, "validate_" + testRootDir, func() {
		vNetSubnets := azure.GetVirtualNetworkSubnets(t, testData.vNetName, testData.vNetRgName, subscriptionID)
		// Ensure subnet is present in the virtual network
		assert.NotNil(t, vNetSubnets[testData.expectedSubnetName])
		// Ensure subnet has the correct address space
		assert.Equal(t, testData.subnetCidr, vNetSubnets[testData.expectedSubnetName])
		// Get the subnet and store in object
		deployedSubnet, err := azure.GetSubnetE(testData.expectedSubnetName, testData.vNetName, testData.vNetRgName, subscriptionID)
		assert.Nil(t, err)
		// Get the subnet's properties
		deployedSubnetProperties := deployedSubnet.SubnetPropertiesFormat
		// Ensure that private endpoint link, and endpoint policies are enabled
		assert.Equal(t, *deployedSubnetProperties.PrivateEndpointNetworkPolicies, "Disabled")
		assert.Equal(t, *deployedSubnetProperties.PrivateLinkServiceNetworkPolicies, "Disabled")
		deployedServiceEndpoints := *deployedSubnetProperties.ServiceEndpoints
		deployedServiceEndpointNames := []string{}
		// Get all service endpoints from deployed subnet
		for i := 0; i < len(deployedServiceEndpoints); i++ {
			deployedServiceEndpointNames = append(deployedServiceEndpointNames, *deployedServiceEndpoints[i].Service)
		}
		// Ensure that all service endpoints are deployed
		for i := 0; i < len(testData.serviceEndpoints); i++ {
			assert.Contains(t, deployedServiceEndpointNames, testData.serviceEndpoints[i])
		}
	})
}

// Loads the nameSuffix (generating one for setup stage if it does not exist)
func GetNameSuffix(t *testing.T, testRootDir string) string {
	nameSuffix := ""
	if (os.Getenv("SKIP_setup_" + testRootDir) == "true") {
		if(ts.IsTestDataPresent(t, testRootDir + ".test-data/nameSuffix.json")) {
			nameSuffix = ts.LoadString(t, testRootDir, "nameSuffix")
		}
	} else {
		nameSuffix = strings.ToLower(randstr.String(suffixLength))
	}	
	return nameSuffix
}

// Copy the module excluding the /test folder and state files
func CopyTerraformFolder(src string, dest string) {
	opt := cp.Options{
		Skip: func(info os.FileInfo, src, dest string) (bool, error) {
			return strings.HasSuffix(src, "/test") || f.PathContainsTerraformState(src), nil
		},
	}
	cp.Copy(src, dest, opt)
}

func TearDown(t *testing.T, testRootDir string) {
	TearDownTerraformOptions(t, testRootDir, testModuleTerraformOptionsDir)
	TearDownTerraformOptions(t, testRootDir, testSetupTerraformOptionsDir)
	os.RemoveAll(testRootDir)
}

func TearDownTerraformOptions(t *testing.T, testRootDir string, terraformOptionsDir string) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from destroy %s: %v\n", terraformOptionsDir, r)
        }
    }()

    if _, err := os.Stat(testRootDir + terraformOptionsDir); err == nil {
        terraformOptions := ts.LoadTerraformOptions(t, fmt.Sprintf("%s%s", testRootDir, terraformOptionsDir))
        _, err := terraform.DestroyE(t, terraformOptions)
		if err != nil {
			panic(err)
		}
    }
}
