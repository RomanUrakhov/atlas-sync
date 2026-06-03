# No SA key file — attached identity via ADC for BigQuery

`bigquery.NewClient` is called with no explicit credentials, which means it resolves auth via Application Default Credentials. In Cloud Run, the GCP metadata server automatically supplies short-lived tokens for the attached service account (`atlas-sync-job`), satisfying ADC with zero code change. We deliberately chose not to use a JSON key file, even though the `.env.example` lists `GOOGLE_APPLICATION_CREDENTIALS` and the developer's existing work pattern involves fetching a key from Secret Manager. Attached identity is strictly better: no long-lived credential to rotate or leak, tokens are auto-rotated by GCP, and the key file pattern is only necessary when you need to act as a different SA cross-project — which is not the case here.

## Considered options

- **JSON key file via `GOOGLE_APPLICATION_CREDENTIALS`** — familiar, works locally and in CI, but introduces a long-lived secret that must be stored, rotated, and kept out of version control.
- **Fetch key from Secret Manager at runtime** — avoids storing the key on disk but still requires a long-lived key to exist somewhere; adds code complexity for no benefit when the compute identity can be attached directly.
- **Attached service account (chosen)** — no key created, no secret to manage, tokens are short-lived and automatic.
