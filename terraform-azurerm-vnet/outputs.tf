output "vnet_name" {
  value       = azurerm_virtual_network.vnet.name
  description = "Vnet name"
}

output "vnet_cidr" {
  value       = azurerm_virtual_network.vnet.address_space
  description = "Vnet CIDR"
}