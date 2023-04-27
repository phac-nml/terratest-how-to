# General configs

variable "location" {
  type        = string
  description = "Azure region for resource deployment. Defaults to canadacentral"
}

variable "ddos_id" {
    type = string
    description = "Distributed denial-of-service plan ID. The plan is located in NMLGC-Core subscription."
}

variable "log_analytics_id" {
    type = string
    description = "Log Analytics workspace ID. The plan is located in NMLGC-Core subscription."
}

# Networking configs

variable "resource_group_name" {
  type        = string
  description = "Resource group that the virtual network lies in"
}

variable "vnet_cidr" {
  type        = list(string)
  description = "The CIDR block definition of the vnet"
}