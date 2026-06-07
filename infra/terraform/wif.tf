# Pool: the namespace of external (non-GCP) identities allowed into this project.
# Issuer-neutral on purpose — a GitLab/AWS provider could join it later without a rename.
resource "google_iam_workload_identity_pool" "ci" {
  workload_identity_pool_id = "cicd"
  display_name              = "CI/CD external identities"
  description               = "External CI/CD identities federated into GCP"
}

# Provider: one trusted token issuer inside the pool above — here, GitHub Actions OIDC.
resource "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.ci.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-oidc"
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.repository" = "assertion.repository"
  }
  attribute_condition = "assertion.repository == '${var.github_repository}'"
  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_service_account_iam_member" "ci_deployer_wif" {
  service_account_id = google_service_account.ci_deployer.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.ci.name}/attribute.repository/${var.github_repository}"
}
