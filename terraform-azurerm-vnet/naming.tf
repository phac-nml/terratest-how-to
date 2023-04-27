locals {
  # Naming locals/constants
  name_prefix = lower(var.name_prefix)
  name_suffix = lower(var.name_suffix)
  vnet_slug   = "vnet"
  vnet_name   = coalesce(var.custom_vnet_name, azurecaf_name.vnet.result)
}

resource "azurecaf_name" "vnet" {
  name          = var.stack
  resource_type = "azurerm_virtual_network"
  prefixes      = var.name_prefix == "" ? null : [local.name_prefix]
  suffixes      = compact([var.client_name, var.environment, local.name_suffix, var.use_caf_naming ? "" : local.vnet_slug])
  use_slug      = var.use_caf_naming
  clean_input   = true
  separator     = "-"
}