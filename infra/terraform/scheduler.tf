resource "google_cloud_scheduler_job" "atlas_sync_trigger" {
  name      = "atlas-sync-trigger"
  region    = var.region
  schedule  = var.scheduler_cron
  time_zone = var.scheduler_timezone
  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project_id}/jobs/${google_cloud_run_v2_job.atlas_sync.name}:run"
    oauth_token {
      service_account_email = google_service_account.job.email
    }
  }
}
