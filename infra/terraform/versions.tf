terraform {
  required_version = ">= 1.6"

  # Remote state lives in the GCS bucket created by infra/bootstrap.sh.
  # The backend block cannot interpolate variables; if your project differs
  # from the default, override at init time:
  #   terraform init -backend-config="bucket=<PROJECT_ID>-tf-state"
  backend "gcs" {
    bucket = "automations-498612-tf-state"
    prefix = "atlas-sync"
  }

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}
