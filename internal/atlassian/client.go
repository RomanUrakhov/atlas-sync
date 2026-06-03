package atlassian

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/RomanUrakhov/atlas-sync/internal/models"
)

type Client struct {
	cloudID    string
	email      string
	apiToken   string
	baseURL    string
	httpClient *http.Client
}

func NewClient(cloudID, email, token, baseURL string) *Client {
	return &Client{
		cloudID:    cloudID,
		email:      email,
		apiToken:   token,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type gqlError struct {
	Message string `json:"message"`
}

type gqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []gqlError      `json:"errors"`
}

// executeQuery sends a GraphQL POST and returns the contents of the "data"
// field. It returns an error for non-200 HTTP responses and for GraphQL-level
// errors (HTTP 200 with an "errors" array in the envelope).
func (c *Client) executeQuery(ctx context.Context, query string, variables map[string]any) (json.RawMessage, error) {
	payload := map[string]any{
		"query":     query,
		"variables": variables,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status %d: %s", resp.StatusCode, body)
	}

	var envelope gqlResponse
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	if len(envelope.Errors) > 0 {
		msgs := make([]string, len(envelope.Errors))
		for i, e := range envelope.Errors {
			msgs[i] = e.Message
		}
		return nil, fmt.Errorf("graphql: %s", strings.Join(msgs, "; "))
	}
	return envelope.Data, nil
}

func (c *Client) containerID() string {
	return "ari:cloud:townsquare::site/" + c.cloudID
}

func (c *Client) TestConnection(ctx context.Context) (json.RawMessage, error) {
	return c.executeQuery(ctx, `query TestConnection { __typename }`, nil)
}

// TODO: potentially rewrite with generics fetchPaginated or smth
func (c *Client) FetchGoals(ctx context.Context) ([]models.Goal, error) {
	var all []models.Goal
	var cursor *string
	for {
		data, err := c.executeQuery(
			ctx,
			`query FetchGoals($containerId: ID!, $first: Int!, $after: String) {
  goals_search(
    containerId: $containerId
    first: $first
    after: $after
  ) {
    edges {
      node {
        id
        name
        status { value }
        targetDate { label }
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}`,
			map[string]any{
				"containerId": c.containerID(),
				"first":       50,
				"after":       cursor,
			},
		)
		if err != nil {
			return nil, err
		}

		var result models.GoalsData
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		for _, edge := range result.GoalSearch.Edges {
			all = append(all, edge.Node)
		}
		if !result.GoalSearch.PageInfo.HasNextPage {
			break
		}
		cursor = &result.GoalSearch.PageInfo.EndCursor
	}
	return all, nil
}

func (c *Client) FetchProjects(ctx context.Context) ([]models.Project, error) {
	var all []models.Project
	var cursor *string
	for {
		data, err := c.executeQuery(
			ctx,
			`query FetchProjects($containerId: String!, $searchString: String!, $first: Int!, $after: String) {
  projects_search(
    containerId: $containerId
	searchString: $searchString
    first: $first
    after: $after
  ) {
    edges {
      node {
        id
        name
        archived
        description { what why }
        dueDate { label }
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}`,
			map[string]any{
				"containerId":  c.containerID(),
				"searchString": "",
				"first":        50,
				"after":        cursor,
			},
		)
		if err != nil {
			return nil, err
		}

		var result models.ProjectsData
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		for _, edge := range result.ProjectSearch.Edges {
			all = append(all, edge.Node)
		}
		if !result.ProjectSearch.PageInfo.HasNextPage {
			break
		}
		cursor = &result.ProjectSearch.PageInfo.EndCursor
	}
	return all, nil
}
