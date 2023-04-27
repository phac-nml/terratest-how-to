## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3.2 |
| <a name="requirement_azurecaf"></a> [azurecaf](#requirement\_azurecaf) | >= 2.30 |
| <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) | >= 3.30 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_azurecaf"></a> [azurecaf](#provider\_azurecaf) | >= 2.30 |
| <a name="provider_azurerm"></a> [azurerm](#provider\_azurerm) | >= 3.30 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [azurecaf_name.vnet](https://registry.terraform.io/providers/aztfmod/azurecaf/latest/docs/resources/name) | resource |
| [azurerm_monitor_diagnostic_setting.diag-vnet](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/monitor_diagnostic_setting) | resource |
| [azurerm_virtual_network.vnet](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/virtual_network) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_client_name"></a> [client\_name](#input\_client\_name) | Client name/account used in naming | `string` | n/a | yes |
| <a name="input_custom_vnet_name"></a> [custom\_vnet\_name](#input\_custom\_vnet\_name) | Optional custom resource vnet name | `string` | `""` | no |
| <a name="input_ddos_id"></a> [ddos\_id](#input\_ddos\_id) | Distributed denial-of-service plan ID. The plan is located in NMLGC-Core subscription. | `string` | n/a | yes |
| <a name="input_environment"></a> [environment](#input\_environment) | Project environment | `string` | n/a | yes |
| <a name="input_location"></a> [location](#input\_location) | Azure region for resource deployment. Defaults to canadacentral | `string` | n/a | yes |
| <a name="input_log_analytics_id"></a> [log\_analytics\_id](#input\_log\_analytics\_id) | Log Analytics workspace ID. The plan is located in NMLGC-Core subscription. | `string` | n/a | yes |
| <a name="input_name_prefix"></a> [name\_prefix](#input\_name\_prefix) | Optional prefix for the generated name | `string` | `""` | no |
| <a name="input_name_suffix"></a> [name\_suffix](#input\_name\_suffix) | Optional suffix for the generated name | `string` | `""` | no |
| <a name="input_resource_group_name"></a> [resource\_group\_name](#input\_resource\_group\_name) | Resource group that the virtual network lies in | `string` | n/a | yes |
| <a name="input_stack"></a> [stack](#input\_stack) | Project stack name | `string` | n/a | yes |
| <a name="input_use_caf_naming"></a> [use\_caf\_naming](#input\_use\_caf\_naming) | Use the Azure CAF naming provider to generate default resource name. | `bool` | `true` | no |
| <a name="input_vnet_cidr"></a> [vnet\_cidr](#input\_vnet\_cidr) | The CIDR block definition of the vnet | `list(string)` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_vnet_cidr"></a> [vnet\_cidr](#output\_vnet\_cidr) | Vnet CIDR |
| <a name="output_vnet_name"></a> [vnet\_name](#output\_vnet\_name) | Vnet name |
