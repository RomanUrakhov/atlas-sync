resource "google_bigquery_dataset" "atlas_sync" {
  dataset_id  = var.dataset_id
  location    = "US"
  description = "atlas-sync: Atlassian Goals & Projects"
}

resource "google_bigquery_table" "goals" {
  dataset_id          = google_bigquery_dataset.atlas_sync.dataset_id
  table_id            = "goals"
  deletion_protection = false
  schema = jsonencode([
    { name = "id", type = "STRING", mode = "REQUIRED" },
    { name = "name", type = "STRING", mode = "NULLABLE" },
    { name = "status", type = "STRING", mode = "NULLABLE" },
    { name = "target_date", type = "STRING", mode = "NULLABLE" }
  ])
}

resource "google_bigquery_table" "projects" {
  dataset_id          = google_bigquery_dataset.atlas_sync.dataset_id
  table_id            = "projects"
  deletion_protection = false
  schema = jsonencode([
    { name = "id", type = "STRING", mode = "REQUIRED" },
    { name = "name", type = "STRING", mode = "NULLABLE" },
    { name = "description_what", type = "STRING", mode = "NULLABLE" },
    { name = "description_why", type = "STRING", mode = "NULLABLE" },
    { name = "due_date", type = "STRING", mode = "NULLABLE" },
    { name = "archived", type = "BOOL", mode = "REQUIRED" }
  ])
}

resource "google_bigquery_dataset_iam_member" "job_data_editor" {
  dataset_id = google_bigquery_dataset.atlas_sync.dataset_id
  role       = "roles/bigquery.dataEditor"
  member     = "serviceAccount:${google_service_account.job.email}"
}

resource "google_project_iam_member" "job_bq_job_user" {
  project = var.project_id
  role    = "roles/bigquery.jobUser"
  member  = "serviceAccount:${google_service_account.job.email}"
}
