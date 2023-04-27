# General naming variables

variable "name_prefix" {
  description = "Optional prefix for the generated name"
  type        = string
  default     = ""
}

variable "name_suffix" {
  description = "Optional suffix for the generated name"
  type        = string
  default     = ""
}

variable "stack" {
  description = "Project stack name"
  type        = string
}

variable "environment" {
  description = "Project environment"
  type        = string
}

variable "client_name" {
  description = "Client name/account used in naming"
  type        = string
}

variable "use_caf_naming" {
  description = "Use the Azure CAF naming provider to generate default resource name. `custom_rg_name` override this if set. Legacy default name is used if this is set to `false`."
  type        = bool
  default     = true
}

# Specific module naming variables

variable "custom_subnet_name" {
  description = "Optional custom resource subnet name"
  type        = string
  default     = ""
}
