resource "google_artifact_registry_repository" "containers" {
  location      = var.region
  repository_id = "atlas-sync"
  format        = "DOCKER"
  description   = "Container images for the atlas-sync Cloud Run Job"
}
