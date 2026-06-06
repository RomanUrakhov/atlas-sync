# Shared `automations` GCP project for all automation flows

atlas-sync is the first of many automation flows we expect to build. Rather than reusing `organic-gecko-329414` (a project created for an unrelated purpose) or creating one GCP project per automation, all automation jobs/services live in a single dedicated `automations` project. The GCP project is our billing, IAM, quota, and Terraform-state boundary, and concentrating automations there means the bootstrap cost — billing link, enabled APIs, Terraform state bucket, and the WIF setup from ADR-0002 — is paid exactly once instead of per flow.

Isolation between flows is achieved softly rather than via separate projects: one service account per automation (least privilege), a dedicated BigQuery dataset per flow, and resource labels for cost attribution. This keeps the strong-isolation benefits we actually need (scoped credentials, per-flow data separation, spend visibility) without re-paying the bootstrap tax. It is also consistent with ADR-0002, which already assumes "all future services reuse the same WIF pool, provider, and SA."

Repository layout (monorepo vs repo-per-flow) is deliberately left open — it is independent of the GCP project boundary and cheap to change later, whereas the project choice is expensive to undo and gates bootstrap.

## Considered options

- **Reuse `organic-gecko-329414`** — zero migration, but atlas-sync squats in a project built for something else, mixing unrelated resources and muddying cost/IAM boundaries.
- **One GCP project per automation** — hard isolation and free per-project cost breakdown, but the full bootstrap tax (billing, APIs, WIF, state bucket) recurs for every flow and invites project sprawl.
- **Shared `automations` project (chosen)** — bootstrap once; soft-isolate per flow via service accounts, BigQuery datasets, and labels. ~80% of the isolation benefit at a fraction of the recurring setup cost.
