locals {
  name_prefix               = var.name_prefix 
  name_suffix               = var.name_suffix
  snet_slug                 = "snet"
  subnet_name               = coalesce(var.custom_subnet_name, azurecaf_name.subnet.result)
  network_security_group_rg = coalesce(var.network_security_group_rg, var.vnet_resource_group_name)
  route_table_rg            = coalesce(var.route_table_rg, var.vnet_resource_group_name)
  network_security_group_id = format("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s", data.azurerm_subscription.current.subscription_id, local.network_security_group_rg, coalesce(var.network_security_group_name, "fake"))
  route_table_id            = format("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/routeTables/%s", data.azurerm_subscription.current.subscription_id, local.route_table_rg, coalesce(var.route_table_name, "fake"))
}

resource "azurecaf_name" "subnet" {
  name          = var.stack
  resource_type = "azurerm_subnet"
  prefixes      = var.name_prefix == "" ? null : [local.name_prefix]
  suffixes      = compact([var.client_name, var.environment, local.name_suffix, var.use_caf_naming ? "" : local.snet_slug])
  use_slug      = var.use_caf_naming
  clean_input   = true
  separator     = "-"
}
