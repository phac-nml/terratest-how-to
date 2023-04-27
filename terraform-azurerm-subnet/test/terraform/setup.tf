terraform {
  backend "local" {}
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = lookup(var.config, "resource_group_name")
  location = lookup(var.config, "location")
}

resource "azurerm_virtual_network" "vnet" {
  name                = lookup(var.config, "vnet_name")
  location            = lookup(var.config, "location")
  address_space       = lookup(var.config, "address_space")
  resource_group_name = azurerm_resource_group.rg.name
}
