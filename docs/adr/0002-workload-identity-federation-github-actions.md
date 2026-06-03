# Workload Identity Federation for GitHub Actions → GCP auth

GitHub Actions CI needs to push Docker images to Artifact Registry and trigger Cloud Run Job deploys. We use Workload Identity Federation (WIF) rather than storing a long-lived SA JSON key in GitHub Secrets. With WIF, GitHub's OIDC provider issues a short-lived token per workflow run; GCP exchanges it for a scoped access token via a configured identity pool and provider. No long-lived credential exists anywhere. This is consistent with the same no-long-lived-keys principle established in ADR-0001 for the runtime identity.

A shared `ci-deployer` service account holds `roles/artifactregistry.writer` and `roles/run.developer` at the project level. All future services reuse the same WIF pool, provider, and SA — the one-time setup cost is paid once.

## Considered options

- **SA JSON key in GitHub Secrets** — simpler to set up, familiar, but creates a long-lived credential that must be rotated and could be exfiltrated from GitHub's secret store.
- **Workload Identity Federation (chosen)** — more complex initial setup (WIF pool, provider, IAM binding), but keyless: tokens are ephemeral, scoped to the workflow run, and require no rotation.
