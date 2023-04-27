<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3.2 |
| <a name="requirement_azurecaf"></a> [azurecaf](#requirement\_azurecaf) | >= 1.2.22 |
| <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) | >= 3.30 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_azurecaf"></a> [azurecaf](#provider\_azurecaf) | >= 1.2.22 |
| <a name="provider_azurerm"></a> [azurerm](#provider\_azurerm) | >= 3.30 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [azurecaf_name.subnet](https://registry.terraform.io/providers/aztfmod/azurecaf/latest/docs/resources/name) | resource |
| [azurerm_subnet.subnet](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subnet) | resource |
| [azurerm_subnet_network_security_group_association.subnet_association](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subnet_network_security_group_association) | resource |
| [azurerm_subnet_route_table_association.route_table_association](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/subnet_route_table_association) | resource |
| [azurerm_subscription.current](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/subscription) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_client_name"></a> [client\_name](#input\_client\_name) | Client name/account used in naming | `string` | n/a | yes |
| <a name="input_custom_subnet_name"></a> [custom\_subnet\_name](#input\_custom\_subnet\_name) | Optional custom resource subnet name | `string` | `""` | no |
| <a name="input_environment"></a> [environment](#input\_environment) | Project environment | `string` | n/a | yes |
| <a name="input_location_short"></a> [location\_short](#input\_location\_short) | Short string for Azure location. | `string` | n/a | yes |
| <a name="input_name_prefix"></a> [name\_prefix](#input\_name\_prefix) | Optional prefix for the generated name | `string` | `""` | no |
| <a name="input_name_suffix"></a> [name\_suffix](#input\_name\_suffix) | Optional suffix for the generated name | `string` | `""` | no |
| <a name="input_network_security_group_name"></a> [network\_security\_group\_name](#input\_network\_security\_group\_name) | The Network Security Group name to associate with the subnets | `string` | `null` | no |
| <a name="input_network_security_group_rg"></a> [network\_security\_group\_rg](#input\_network\_security\_group\_rg) | The Network Security Group RG to associate with the subnet. Default is the same RG than the subnet. | `string` | `null` | no |
| <a name="input_private_endpoint_enabled"></a> [private\_endpoint\_enabled](#input\_private\_endpoint\_enabled) | Enable or Disable network policies for the private endpoint on the subnet. Setting this to true will Enable the policy and setting this to false will Disable the policy. Defaults to true. | `bool` | `true` | no |
| <a name="input_private_link_service_enabled"></a> [private\_link\_service\_enabled](#input\_private\_link\_service\_enabled) | Enable or Disable network policies for the private link service on the subnet. Setting this to true will Enable the policy and setting this to false will Disable the policy. Defaults to true. | `bool` | `true` | no |
| <a name="input_route_table_name"></a> [route\_table\_name](#input\_route\_table\_name) | The Route Table name to associate with the subnet | `string` | `null` | no |
| <a name="input_route_table_rg"></a> [route\_table\_rg](#input\_route\_table\_rg) | The Route Table RG to associate with the subnet. Default is the same RG than the subnet. | `string` | `null` | no |
| <a name="input_service_endpoints"></a> [service\_endpoints](#input\_service\_endpoints) | The list of Service endpoints to associate with the subnet | `list(string)` | `[]` | no |
| <a name="input_stack"></a> [stack](#input\_stack) | Project stack name | `string` | n/a | yes |
| <a name="input_subnet_cidr_list"></a> [subnet\_cidr\_list](#input\_subnet\_cidr\_list) | The address prefix list to use for the subnet | `list(string)` | n/a | yes |
| <a name="input_subnet_delegation"></a> [subnet\_delegation](#input\_subnet\_delegation) | Configuration delegations on subnet<br>object({<br>  name = object({<br>    name = string,<br>    actions = list(string)<br>  })<br>}) | `map(list(any))` | `{}` | no |
| <a name="input_use_caf_naming"></a> [use\_caf\_naming](#input\_use\_caf\_naming) | Use the Azure CAF naming provider to generate default resource name. `custom_rg_name` override this if set. Legacy default name is used if this is set to `false`. | `bool` | `true` | no |
| <a name="input_vnet_name"></a> [vnet\_name](#input\_vnet\_name) | Virtual network name | `string` | n/a | yes |
| <a name="input_vnet_resource_group_name"></a> [vnet\_resource\_group\_name](#input\_vnet\_resource\_group\_name) | Resource group name | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_subnet_cidr_list"></a> [subnet\_cidr\_list](#output\_subnet\_cidr\_list) | CIDR list of the created subnets |
| <a name="output_subnet_cidrs_map"></a> [subnet\_cidrs\_map](#output\_subnet\_cidrs\_map) | Map with names and CIDRs of the created subnets |
| <a name="output_subnet_id"></a> [subnet\_id](#output\_subnet\_id) | Id of the created subnet |
| <a name="output_subnet_ips"></a> [subnet\_ips](#output\_subnet\_ips) | The collection of IPs within this subnet |
| <a name="output_subnet_names"></a> [subnet\_names](#output\_subnet\_names) | Names of the created subnet |
<!-- END_TF_DOCS -->