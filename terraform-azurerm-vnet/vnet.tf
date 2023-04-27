resource "azurerm_virtual_network" "vnet" {
  name                = local.vnet_name
  location            = var.location
  resource_group_name = var.resource_group_name
  address_space       = var.vnet_cidr
  
  ddos_protection_plan {
    id     = var.ddos_id
    enable = "true"
  } 
}

resource "azurerm_monitor_diagnostic_setting" "diag-vnet" {
  name               = "diag-${local.vnet_name}"
  target_resource_id = azurerm_virtual_network.vnet.id
  log_analytics_workspace_id = var.log_analytics_id

  log {
    category = "VMProtectionAlerts"
    
    retention_policy {
      enabled = "false"
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = "false"
    }
  }
}