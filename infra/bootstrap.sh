#!/usr/bin/env bash
#
# One-time GCP bootstrap for atlas-sync.
#
# Creates the prerequisites that Terraform itself cannot manage:
#   1. Enables the GCP APIs the deployment depends on.
#   2. Creates the GCS bucket that holds Terraform remote state (versioned).
#
# This is deliberately NOT Terraform — you cannot Terraform the backend that
# Terraform stores its state in. Run this once per project, by hand, with an
# authenticated gcloud (`gcloud auth login` + `gcloud config set project ...`).
#
# The script is idempotent: enabling an already-enabled API and re-running
# against an existing bucket are both no-ops.
#
# Usage:
#   ./infra/bootstrap.sh
#   PROJECT_ID=my-project REGION=europe-west1 ./infra/bootstrap.sh

set -euo pipefail

PROJECT_ID="${PROJECT_ID:-automations-498612}"
REGION="${REGION:-us-central1}"
# State bucket names must be globally unique; project-id prefix guarantees it.
STATE_BUCKET="${STATE_BUCKET:-${PROJECT_ID}-tf-state}"

# APIs required by the deployment. Keep this list in sync with infra/terraform.
APIS=(
  run.googleapis.com              # Cloud Run Job (the sync runtime)
  cloudscheduler.googleapis.com   # scheduled trigger for the job
  secretmanager.googleapis.com    # Atlassian credentials at runtime
  artifactregistry.googleapis.com # container image registry
  iam.googleapis.com              # service accounts + IAM bindings
  bigquery.googleapis.com         # destination dataset/tables
  # Needed by Workload Identity Federation for GitHub Actions (issue #5):
  iamcredentials.googleapis.com   # token exchange for federated identity
  sts.googleapis.com              # Security Token Service for WIF
)

echo "==> Project: ${PROJECT_ID}"
echo "==> Region:  ${REGION}"
echo "==> State bucket: gs://${STATE_BUCKET}"
echo

gcloud config set project "${PROJECT_ID}" >/dev/null

echo "==> Enabling APIs (idempotent)..."
gcloud services enable "${APIS[@]}" --project "${PROJECT_ID}"
echo "    done."
echo

echo "==> Ensuring Terraform state bucket exists..."
if gcloud storage buckets describe "gs://${STATE_BUCKET}" >/dev/null 2>&1; then
  echo "    gs://${STATE_BUCKET} already exists — skipping create."
else
  gcloud storage buckets create "gs://${STATE_BUCKET}" \
    --project "${PROJECT_ID}" \
    --location "${REGION}" \
    --uniform-bucket-level-access
  echo "    created gs://${STATE_BUCKET}."
fi

# Versioning lets us recover a corrupted/clobbered state file.
echo "==> Enabling object versioning on the state bucket..."
gcloud storage buckets update "gs://${STATE_BUCKET}" --versioning
echo "    done."
echo

echo "==> Bootstrap complete. Next:"
echo "      cd infra/terraform && terraform init"
