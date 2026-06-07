resource "google_service_account" "job" {
  account_id   = "atlas-sync-job"
  display_name = "atlas-sync Cloud Run Job runtime"
}

resource "google_service_account" "ci_deployer" {
  account_id   = "ci-deployer"
  display_name = "CI/CD deployer for GitHub actions"
}

resource "google_project_iam_member" "ci_deployer_artifactregistry" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.ci_deployer.email}"
}

resource "google_project_iam_member" "ci_deployer_run" {
  project = var.project_id
  role    = "roles/run.developer"
  member  = "serviceAccount:${google_service_account.ci_deployer.email}"
}

resource "google_service_account_iam_member" "ci_deployer_act_as_job" {
  service_account_id = google_service_account.job.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.ci_deployer.email}"
}
