provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = lookup(var.config, "resource_group_name")
  location = lookup(var.config, "location")
}
