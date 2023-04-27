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

// Global test constants
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
	testModuleTerraformOptionsDir = "virtualNetworkTerraformOptions/"
	testSetupTerraformOptionsDir = "setupTerraformOptions/"
	ddosPlanID         = "<private_value>"
	laWorkspaceID      = "<private_value>"
)

// A struct containing any variables needed for implementing a test
type VirtualNetworkTestData struct {
	vNetRgName string
	vNetCidr []string
	vNetName string
	vNetDDOSID string
	vNetLAWorkspaceID string
}

func TestVirtualNetwork(t *testing.T) {
	testRootDir := "TestVirtualNetwork/"

	// Uncomment any of the following lines to skip that test stage
	// os.Setenv("SKIP_setup_" + testRootDir, "true")
	// os.Setenv("SKIP_deploy_" + testRootDir, "true")
	// os.Setenv("SKIP_validate_" + testRootDir, "true")
	// os.Setenv("SKIP_teardown_" + testRootDir, "true")

	t.Parallel() // Remove to test serially

	nameSuffix := GetNameSuffix(t, testRootDir)

	testData := VirtualNetworkTestData {
		vNetRgName: fmt.Sprintf("rg-vnet-unit-test-%s", nameSuffix),
		vNetCidr: []string{"10.0.0.0/16"},
		vNetName: fmt.Sprintf("vnet-stack-client-test-%s", nameSuffix),
		vNetDDOSID: ddosPlanID,
		vNetLAWorkspaceID: laWorkspaceID,
	}

	VirtualNetwork(t, testRootDir, nameSuffix, testData)
}

func VirtualNetwork(t *testing.T, testRootDir string, nameSuffix string, testData VirtualNetworkTestData) {
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
				},
			},
		})

		ts.SaveTerraformOptions(t, fmt.Sprintf("%s%s", testRootDir, testSetupTerraformOptionsDir), setupTerraformOptions)
		terraform.InitAndApply(t, setupTerraformOptions)
	})

	ts.RunTestStage(t, "deploy_" + testRootDir, func() {
		CopyTerraformFolder(moduleTerraformDir, fmt.Sprintf("%s%s", testRootDir, testModuleDir))

		virtualNetworkTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: fmt.Sprintf("%s%s", testRootDir, testModuleDir),
			Vars: map[string]interface{}{
				"location":            location,
				"ddos_id":             testData.vNetDDOSID,
				"log_analytics_id":    testData.vNetLAWorkspaceID,
				"vnet_cidr":           testData.vNetCidr,
				"resource_group_name": testData.vNetRgName,
				"name_suffix":         nameSuffix,
				"client_name":         clientName,
				"environment":         environment,
				"stack":               stack,
			},
		})

		ts.SaveTerraformOptions(t, fmt.Sprintf("%s%s", testRootDir, testModuleTerraformOptionsDir), virtualNetworkTerraformOptions)
		terraform.InitAndApply(t, virtualNetworkTerraformOptions)
	})

	ts.RunTestStage(t, "validate_" + testRootDir, func() {
		// Assert that the virtual network exists
		assert.True(t, azure.VirtualNetworkExists(t, testData.vNetName, testData.vNetRgName, subscriptionID))

		// Get the deployed virtual network properties
		deployedVNet, err := azure.GetVirtualNetworkE(testData.vNetName, testData.vNetRgName, subscriptionID)

		// Basic assertions to ensure no errors, and proper attributes are correct
		assert.Nil(t, err)
		assert.NotNil(t, *deployedVNet.ID)
		assert.Equal(t, testData.vNetName, *deployedVNet.Name)
		assert.Equal(t, strings.ToLower(location), strings.ToLower(*deployedVNet.Location))

		// Virtual network address configs
		deployedVNetAddrConfs := deployedVNet.VirtualNetworkPropertiesFormat
		assert.Equal(t, testData.vNetCidr, *deployedVNetAddrConfs.AddressSpace.AddressPrefixes)
		assert.True(t, *deployedVNetAddrConfs.EnableDdosProtection)
		assert.Equal(t, testData.vNetDDOSID, *deployedVNetAddrConfs.DdosProtectionPlan.ID)
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