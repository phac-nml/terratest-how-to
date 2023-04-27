# terratest-how-to

This repository showcases the utilization of Terratest for testing Terraform modules. Given the scarcity of examples illustrating testing implementation on Terraform modules, the purpose of this repository is to provide testing examples taken from our internal repositories. Specifically, it provides two examples for testing Azure resources: [azurerm_virtual_network](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/virtual_network) and [azurerm_subnet](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subnet).

# Terratest

Terratest is a Go testing framework that can be used for testing Terraform modules. It is definetely worth watching this [video](https://www.youtube.com/watch?v=xhHOW0EF5u8) and taking a look at the [documentation](https://terratest.gruntwork.io/docs/#getting-started) before jumping into Terratest. 

# Getting Started

The following dependencies are required:

- [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli)
- [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
- [Golang](https://go.dev/doc/install)

Login and set the correct Azure Subscription using the Azure CLI:
```
az login
az account set --subscription <subscription_id_or_name>
```

# Implementing Terratest on an Existing Module

1. Create the `test` folder at the root of the module directory.
2. Create a new file called `terraform-name-of-module_test.go`. (eg. `terraform-azurerm-vnet_test.go`)
3. Copy the code from `example_test.go` to your new `_test.go` file.
4. Configure Go dependencies by running:
```
cd test
go mod init terraform-name-of-module
go mod tidy
```
5. Begin writing tests!

# Terraform Module Structure

The resulting Terraform module folder structure should look like this:
```
/terraform-azurerm-<resource_name>
├── main.tf
├── naming.tf
├── outputs.tf
├── provider.tf
├── README.md
├── test <------------ This is the directory that contains the Terratest file (_test.go)
│   ├── go.mod
│   ├── go.sum
│   ├── terraform <--- This is the directory that contains pre-requisite Terraform resources for the test
│   │   ├── outputs.tf
│   │   ├── setup.tf
│   │   └── variables.tf
│   └── terraform-azurerm-storage-account_test.go <-- The actual Terratest file itself
├── variables-naming.tf
├── variables.tf
└── versions.tf
```

# Testing Pattern

The following testing stages should be used to test each Terraform module. These stages can be run independently of each other.

1. **Setup:** deploy resource dependencies
   - This is not necessary if testing an "all-in-one" module (i.e. there are no setup resources)
2. **Deploy:** deploy the module infrastructure that is being tested
3. **Validate:** test the deployed infrastructure to ensure that it works correctly
4. **Teardown:** destroy (undeploy) all infrastructure and remove all Terraform state files

# Test Structure

Every `terraform-name-of-module_test.go` test file will contain the following sections:

1. Imports
2. Global Test Constants
3. Test Data Struct(s)
4. Test Function(s)
5. Test Helper Function(s)
6. Private functions used across all tests (eg. `GetNameSuffix()`)

# Test Workflow

## Test Data Struct(s)

A `testData` struct is used to neatly contain all variables used for testing. Since test stages can be run independently, it is useful to pass test variables as a parameter, rather than storing test variables in a file.

## Test Function

The test function must start with "Test" for the test to be called when running `go test`. Each test function is structured as follows:

```
func Test<name_of_test>(t *testing.T) {
    1. Set `testRootDir` (all testing occurs here)
    2. Set env variables to skip test stages (for quick local testing)
    3. Set tests to run in parallel (can be commented out for local testing)
    4. Generate a random nameSuffix to avoid resource naming collisions (or load existing nameSuffix)
    5. Initialize testData using a struct
    6. Call <name_of_test>()
}
```

## Test Helper Function

The test helper function is called by the test function. This function runs whichever test stages are not set to "skipped" via environment variables:

```
func <name_of_test>(t *testing.T, testRootDir string, nameSuffix string, testData TestData) {
    1. Run defer teardown stage
    2. Run setup stage
    3. Run deploy stage
    4. Run validate stage
}
```

## Test Stages

### Setup

1. Run teardown to reinitialize setup (destroys any existing resources and removes test folders containing state)
2. Copy `terraform` setup folder to `testRootDir`
3. Save `nameSuffix` to use in later test runs (if setup is skipped)
4. Create and save `setupTerraformOptions`
5. Initialize and apply `setupTerraformOptions`

### Deploy

1. Copy the module (eg. `terraform-azurerm-vnet`) folder excluding `/test` to `testRootDir`
2. Create and save `moduleTerraformOptions`
   - It is sometimes necessary to load `setupTerraformOptions` and use `terraform.Output()` to access dynamic variables created in `setup` that are needed for `deploy`
3. Initialize and apply `moduleTerraformOptions`

### Validate

1. Use assertions to validate deployed infrastructure
2. Optionally add `GetAzureResourceClient()` to `terraform-name-of-module_test.go` to test more resource fields including Name, ID, Tags and Properties
    - Can be useful as the built-in Terratest functions are somewhat limited and do not cover all Azure resources
    - The necessary parameters to pass to GetAzureResourceClient() can be found in the JSON view of a deployed resource through the Azure Portal

```
func GetAzureResourceClient(subscriptionID string, resourceGroupName string, resourceProviderNamespace string, parentResourceType string, parentResourceName string, resourceType string, resourceName string, apiVersion string) (*armresources.ClientGetResponse, error){ 
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if(err != nil) {
		return nil, err
	}

	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if(err != nil) {
		return nil, err
	}

	parentResourcePath := ""
	if(parentResourceType != "") {
		parentResourcePath = fmt.Sprintf("%s/%s", parentResourceType, parentResourceName)
	}
	resourceClient, err := client.Get(context.TODO(), resourceGroupName, resourceProviderNamespace, parentResourcePath, resourceType, resourceName, apiVersion, nil)
	if(err != nil) {
		return nil, err
	}

	return &resourceClient, nil
}
```

### Teardown

1. Destroy `moduleTerraformOptions` (recover if non-existent)
2. Destroy `setupTerraformOptions` (recover if non-existent)
3. Remove the `testRootDir` that contains all copied files and Terraform state

## Running Tests

Run all tests in `terraform-name-of-module_test.go` with:

```
cd test
go test -v -timeout 30m | tee test_output.log
terratest_log_parser -testlog test_output.log -outputdir test_output
```

Note that when running tests in parallel it is necessary to parse the interleaved log output as done above. The Terratest Log Parser will create a `report.xml` file that can be used to integrate with CircleCI or Azure DevOps. See more information [here](https://terratest.gruntwork.io/docs/testing-best-practices/debugging-interleaved-test-output/).

## Common Testing Approach
1. Run just the `setup` stage until the setup resources deploy correctly (setup resources will be destroyed on each run)
    - Can also manually delete the resource group through the portal and delete the `testRootDir` for faster iterating
2. Run `setup` and `deploy` until the module deploys correctly
3. Once deployed run `validate` only until satisfied with test assertions (avoids redeploying each time)
    - Running just `validate` should take less than 10 seconds
4. Run all stages to ensure everything works properly