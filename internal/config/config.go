package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AtlassianCloudID  string
	AtlassianEmail    string
	AtlassianAPIToken string
	BigQueryProject   string
	BigQueryDataset   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AtlassianCloudID:  os.Getenv("ATLASSIAN_CLOUD_ID"),
		AtlassianEmail:    os.Getenv("ATLASSIAN_EMAIL"),
		AtlassianAPIToken: os.Getenv("ATLASSIAN_API_TOKEN"),
		BigQueryProject:   os.Getenv("BIGQUERY_PROJECT"),
		BigQueryDataset:   os.Getenv("BIGQUERY_DATASET"),
	}

	for _, check := range []struct{ val, name string }{
		{cfg.AtlassianCloudID, "ATLASSIAN_CLOUD_ID"},
		{cfg.AtlassianEmail, "ATLASSIAN_EMAIL"},
		{cfg.AtlassianAPIToken, "ATLASSIAN_API_TOKEN"},
		{cfg.BigQueryProject, "BIGQUERY_PROJECT"},
		{cfg.BigQueryDataset, "BIGQUERY_DATASET"},
	} {
		if check.val == "" {
			return nil, fmt.Errorf("%s env variable is required", check.name)
		}
	}

	return cfg, nil
}
