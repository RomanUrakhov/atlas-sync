locals {
  atlassian_secret_ids = [
    "atlas-sync-atlassian-api-token",
    "atlas-sync-atlassian-cloud-id",
    "atlas-sync-atlassian-email"
  ]
}


resource "google_secret_manager_secret" "atlassian" {
  for_each  = toset(local.atlassian_secret_ids)
  secret_id = each.value
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_iam_member" "job_accessor" {
  for_each  = google_secret_manager_secret.atlassian
  secret_id = each.value.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.job.email}"
}
