# atlas-sync

A CLI tool that pulls Goals and Projects from an Atlassian cloud workspace and loads them into BigQuery for reporting.

## Language

**Sync**:
A single run-to-completion operation: fetch all Goals and Projects from the Atlassian GraphQL API and load them into BigQuery. The tool does one sync per invocation.
_Avoid_: import, export, transfer, job (use "sync" for the operation; "job" is reserved for the Cloud Run Job infrastructure concept)

**Goal**:
An Atlassian Goals entity — a trackable outcome belonging to a workspace. Represented in BigQuery as a row in the `goals` table.
_Avoid_: objective, OKR

**Project**:
An Atlassian Projects entity — a trackable initiative belonging to a workspace. Represented in BigQuery as a row in the `projects` table.
_Avoid_: initiative, epic

**Cloud ID**:
The unique identifier for an Atlassian cloud workspace instance. Required to scope all GraphQL API requests to the correct workspace.
_Avoid_: workspace ID, tenant ID, org ID
