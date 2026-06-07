variable "project_id" {
  description = "the shared automations project"
  type        = string
  default     = "automations-498612"
}

variable "region" {
  description = "the shared automations project region"
  type        = string
  default     = "us-central1"
}

variable "dataset_id" {
  description = "the BigQuery dataset for atlas-sync automation"
  type        = string
  default     = "atlas_sync"
}

variable "image" {
  description = "placeholder image"
  type        = string
  default     = "us-docker.pkg.dev/cloudrun/container/job"
}

variable "scheduler_cron" {
  description = "daily 6am"
  type        = string
  default     = "0 6 * * *"
}

variable "scheduler_timezone" {
  description = "default cron timezone"
  type        = string
  default     = "Asia/Kuala_Lumpur"
}

variable "github_repository" {
  description = "owner/repo allowed to impersonate ci-deployer via WIF"
  type        = string
  default     = "RomanUrakhov/atlas-sync"
}
