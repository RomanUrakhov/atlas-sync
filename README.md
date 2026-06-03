# atlas-sync

Pulls Goals and Projects from Atlassian's GraphQL API and lands them in BigQuery. Built as a Go learning project — nothing fancy, just a clean sync pipeline.

## How it works

```
┌─────────────────────────────────────────────────────────┐
│                     atlas-sync sync                     │
└───────────────────────┬─────────────────────────────────┘
                        │ load env config
                        ▼
              ┌─────────────────┐
              │  config.Load()  │
              └────────┬────────┘
                       │
          ┌────────────┴────────────┐
          ▼                         ▼
  ┌──────────────┐         ┌────────────────┐
  │  Atlassian   │         │   BigQuery     │
  │   Client     │         │    Writer      │
  │  (GraphQL)   │         │  (load jobs)   │
  └──────┬───────┘         └───────┬────────┘
         │                         │
         └───────────┬─────────────┘
                     ▼
             ┌───────────────┐
             │    Syncer     │
             │               │
             │ FetchGoals    │──► InsertGoals
             │ FetchProjects │──► InsertProjects
             └───────────────┘
```

The `Syncer` is the heart of it — it holds a `Fetcher` interface (the Atlassian client) and a `Writer` interface (BigQuery), so either side can be swapped for testing. Pagination is handled inside the client using cursor-based loops (50 records per page).

## Setup

You'll need a [personal Atlassian API token](https://id.atlassian.com/manage-profile/security/api-tokens) and a BigQuery dataset with a service account key.

```bash
export ATLASSIAN_EMAIL=you@example.com
export ATLASSIAN_API_TOKEN=your-token
export ATLASSIAN_CLOUD_ID=your-cloud-id
export BIGQUERY_PROJECT=your-gcp-project
export BIGQUERY_DATASET=your-dataset
export BIGQUERY_TABLE=your-table
```

## Usage

```bash
make run          # build + sync
make build        # build only
make test         # unit tests
make test-integration  # integration tests (needs live credentials)
```

Or directly:

```bash
./atlas-sync sync
```

## Potential improvements

- Concurrent fetching of goals and projects
- Idempotent writes — dedup against existing BigQuery rows before inserting
- Retry with exponential backoff for transient API/BQ errors
- Scale to additional Atlassian resource types without restructuring the pipeline
- CI/CD pipeline (GitHub Actions)
- Deploy as a Cloud Run Job on GCP, triggered by Cloud Scheduler

## Project structure

```
cmd/atlas-sync/       # CLI entry point (cobra)
internal/
  atlassian/          # GraphQL client, FetchGoals, FetchProjects
  bqwriter/           # BigQuery writer + in-memory fake for tests
  config/             # env var loading
  models/             # typed structs for Atlassian responses
  syncer/             # orchestration logic
```
