//go:build integration

package atlassian_test

import (
	"context"
	"os"
	"testing"

	"github.com/RomanUrakhov/atlas-sync/internal/atlassian"
	"github.com/joho/godotenv"
)

func newClientFromEnv(t *testing.T) *atlassian.Client {
	t.Helper()
	_ = godotenv.Load("../../.env")

	cloudID := os.Getenv("ATLASSIAN_CLOUD_ID")
	email := os.Getenv("ATLASSIAN_EMAIL")
	token := os.Getenv("ATLASSIAN_API_TOKEN")

	if token == "" {
		t.Skip("ATLASSIAN_API_TOKEN not set, skipping integration test")
	}
	return atlassian.NewClient(cloudID, email, token, "https://api.atlassian.com/graphql")
}

func TestFetchGoals_Integration(t *testing.T) {
	client := newClientFromEnv(t)
	goals, err := client.FetchGoals(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("fetched %d goals", len(goals))
}

func TestFetchProjects_Integration(t *testing.T) {
	client := newClientFromEnv(t)
	projects, err := client.FetchProjects(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("fetched %d projects", len(projects))
}
