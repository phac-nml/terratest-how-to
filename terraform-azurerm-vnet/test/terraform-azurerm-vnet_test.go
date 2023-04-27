package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	ts "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

// Global test variables
var (
	subscriptionID     = "68dfddb8-d4ed-4f1c-bd08-78674c7088f7"
	location           = "CanadaCentral"
	clientName         = "client"
	environment        = "test"
	stack              = "stack"
	suffixLength       = 8
	setupTerraformDir  = "terraform"
	moduleTerraformDir = "../"
	ddosPlanID         = "/subscriptions/fdc723e3-0c24-4b90-bb60-55bf4c819bf2/resourceGroups/rg-nmlgc-rz-security/providers/Microsoft.Network/ddosProtectionPlans/ddos-nmlgc-protection-plan"
	laWorkspaceID      = "/subscriptions/fdc723e3-0c24-4b90-bb60-55bf4c819bf2/resourceGroups/rg-nmlgc-rz-security/providers/Microsoft.OperationalInsights/workspaces/workspace-nmlgc-all"
)

func TestVirtualNetworkSingleCIDR(t *testing.T) {
	// Uncomment any of the following lines to skip that test stage
	// os.Setenv("SKIP_setup", "true")
	// os.Setenv("SKIP_vnetApply", "true")
	// os.Setenv("SKIP_validate", "true")
	// os.Setenv("SKIP_teardown", "true")

	// Generate test variables
	nameSuffix := strings.ToLower(randstr.String(suffixLength))
	vNetRgName := fmt.Sprintf("rg-vnet-unit-test-%s", nameSuffix)
	// Variables for vnet test
	vNetCidr := []string{"10.0.0.0/16"}
	vNetName := fmt.Sprintf("vnet-stack-client-test-%s", nameSuffix)
	vNetDDOSID := ddosPlanID
	vNetLAWorkspaceID := laWorkspaceID
	// Terraform options for setup resources
	setupTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: setupTerraformDir,
		Vars: map[string]interface{}{
			"config": map[string]interface{}{
				"location":            location,
				"resource_group_name": vNetRgName,
			},
		},
	})
	// Terraform options for vnet resources
	vnetTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: moduleTerraformDir,
		Vars: map[string]interface{}{
			"location":            location,
			"ddos_id":             vNetDDOSID,
			"log_analytics_id":    vNetLAWorkspaceID,
			"vnet_cidr":           vNetCidr,
			"resource_group_name": vNetRgName,
			"name_suffix":         nameSuffix,
			"client_name":         clientName,
			"environment":         environment,
			"stack":               stack,
		},
	})

	// At the end of the test, clean up resources.
	defer ts.RunTestStage(t, "teardown", func() {
		terraform.Destroy(t, vnetTerraformOptions)
		terraform.Destroy(t, setupTerraformOptions)
	})

	ts.RunTestStage(t, "setup", func() {
		terraform.InitAndApply(t, setupTerraformOptions)
	})

	ts.RunTestStage(t, "vnetApply", func() {
		terraform.InitAndApply(t, vnetTerraformOptions)
	})

	ts.RunTestStage(t, "validate", func() {
		// Assert that the virtual network exists
		assert.True(t, azure.VirtualNetworkExists(t, vNetName, vNetRgName, subscriptionID))

		// Get the deployed virtual network properties
		deployedVNet, err := azure.GetVirtualNetworkE(vNetName, vNetRgName, subscriptionID)

		// Basic assertions to ensure no errors, and proper attributes are correct
		assert.Nil(t, err)
		assert.NotNil(t, *deployedVNet.ID)
		assert.Equal(t, vNetName, *deployedVNet.Name)
		assert.Equal(t, strings.ToLower(location), strings.ToLower(*deployedVNet.Location))

		// Virtual network address configs
		deployedVNetAddrConfs := deployedVNet.VirtualNetworkPropertiesFormat
		assert.Equal(t, vNetCidr, *deployedVNetAddrConfs.AddressSpace.AddressPrefixes)
		assert.True(t, *deployedVNetAddrConfs.EnableDdosProtection)
		assert.Equal(t, vNetDDOSID, *deployedVNetAddrConfs.DdosProtectionPlan.ID)
	})
}

func TestVirtualNetworkMultipleCIDR(t *testing.T) {
	// Uncomment any of the following lines to skip that test stage
	// os.Setenv("SKIP_setup", "true")
	// os.Setenv("SKIP_vnetApply", "true")
	// os.Setenv("SKIP_validate", "true")
	// os.Setenv("SKIP_teardown", "true")

	// Generate test variables
	nameSuffix := strings.ToLower(randstr.String(suffixLength))
	vNetRgName := fmt.Sprintf("rg-vnet-unit-test-%s", nameSuffix)
	// Variables for vnet test
	vNetCidr := []string{"10.5.0.0/24", "10.6.0.0/24"}
	vNetName := fmt.Sprintf("vnet-stack-client-test-%s", nameSuffix)
	vNetDDOSID := ddosPlanID
	vNetLAWorkspaceID := laWorkspaceID
	// Terraform options for setup resources
	setupTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: setupTerraformDir,
		Vars: map[string]interface{}{
			"config": map[string]interface{}{
				"location":            location,
				"resource_group_name": vNetRgName,
			},
		},
	})
	// Terraform options for vnet resources
	vnetTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: moduleTerraformDir,
		Vars: map[string]interface{}{
			"location":            location,
			"ddos_id":             vNetDDOSID,
			"log_analytics_id":    vNetLAWorkspaceID,
			"vnet_cidr":           vNetCidr,
			"resource_group_name": vNetRgName,
			"name_suffix":         nameSuffix,
			"client_name":         clientName,
			"environment":         environment,
			"stack":               stack,
		},
	})

	// At the end of the test, clean up resources.
	defer ts.RunTestStage(t, "teardown", func() {
		terraform.Destroy(t, vnetTerraformOptions)
		terraform.Destroy(t, setupTerraformOptions)
	})

	ts.RunTestStage(t, "setup", func() {
		terraform.InitAndApply(t, setupTerraformOptions)
	})

	ts.RunTestStage(t, "vnetApply", func() {
		terraform.InitAndApply(t, vnetTerraformOptions)
	})

	ts.RunTestStage(t, "validate", func() {
		// Assert that the virtual network exists
		assert.True(t, azure.VirtualNetworkExists(t, vNetName, vNetRgName, subscriptionID))

		// Get the deployed virtual network properties
		deployedVNet, err := azure.GetVirtualNetworkE(vNetName, vNetRgName, subscriptionID)

		// Basic assertions to ensure no errors, and proper attributes are correct
		assert.Nil(t, err)
		assert.NotNil(t, *deployedVNet.ID)
		assert.Equal(t, vNetName, *deployedVNet.Name)
		assert.Equal(t, strings.ToLower(location), strings.ToLower(*deployedVNet.Location))

		// Virtual network address configs
		deployedVNetAddrConfs := deployedVNet.VirtualNetworkPropertiesFormat
		assert.Equal(t, vNetCidr, *deployedVNetAddrConfs.AddressSpace.AddressPrefixes)
		assert.True(t, *deployedVNetAddrConfs.EnableDdosProtection)
		assert.Equal(t, vNetDDOSID, *deployedVNetAddrConfs.DdosProtectionPlan.ID)
	})
}
