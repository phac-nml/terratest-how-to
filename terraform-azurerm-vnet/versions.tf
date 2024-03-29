terraform {
  required_version = ">= 1.3.2"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.30"
    }
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = ">= 1.2.22"
    }  
  }
}