output "runtime_service_account" {
  description = "Email of the atlas-sync-job SA the Cloud Run Job runs as. Seed Atlassian secrets so this identity can read them (issue #6)."
  value       = google_service_account.job.email
}

output "ci_deployer_service_account" {
  description = "Email of the ci-deployer SA used by CI/CD. Bind it to the Workload Identity Federation pool for GitHub Actions (issue #5)."
  value       = google_service_account.ci_deployer.email
}

output "image_repository" {
  description = "Fully-qualified image path the CD pipeline pushes to and the Cloud Run Job pulls from (issue #7). Append a tag, e.g. :<git-sha>."
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.containers.repository_id}/atlas-sync"
}
