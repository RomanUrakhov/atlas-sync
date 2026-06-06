resource "google_cloud_run_v2_job" "atlas_sync" {
  name     = "atlas-sync"
  location = var.region
  template {
    template {
      service_account = google_service_account.job.email
      containers {
        image = var.image
        env {
          name  = "BIGQUERY_PROJECT"
          value = var.project_id
        }
        env {
          name  = "BIGQUERY_DATASET"
          value = var.dataset_id
        }
        env {
          name = "ATLASSIAN_API_TOKEN"
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.atlassian["atlas-sync-atlassian-api-token"].secret_id
              version = "latest"
            }
          }
        }
        env {
          name = "ATLASSIAN_EMAIL"
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.atlassian["atlas-sync-atlassian-email"].secret_id
              version = "latest"
            }
          }
        }
        env {
          name = "ATLASSIAN_CLOUD_ID"
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.atlassian["atlas-sync-atlassian-cloud-id"].secret_id
              version = "latest"
            }
          }
        }
      }
    }

  }
  lifecycle {
    ignore_changes = [
      template[0].template[0].containers[0].image,
    ]
  }
}

resource "google_cloud_run_v2_job_iam_member" "invoker" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_job.atlas_sync.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.job.email}"
}
