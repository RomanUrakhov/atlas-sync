package atlassian_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RomanUrakhov/atlas-sync/internal/atlassian"
)

func TestFetchGoals(t *testing.T) {
	fakeResponse := `{
    "data": {
        "goals_search": {
            "edges": [
                {
                    "node": {
                        "name": "First Go service deployed to GCP",
                        "status": {
                            "value": "pending"
                        },
                        "targetDate": {
                            "label": "August"
                        },
                        "id": "ari:cloud:townsquare:ef524de6-fdee-4937-a28b-5d531545cbae:goal/a1a46457-943c-4df3-bf05-744ab089e271"
                    }
                },
                {
                    "node": {
                        "name": "Monorepo live with all services migrated",
                        "status": {
                            "value": "pending"
                        },
                        "targetDate": {
                            "label": "August"
                        },
                        "id": "ari:cloud:townsquare:ef524de6-fdee-4937-a28b-5d531545cbae:goal/961c7737-2744-4979-9d1e-234e131372d0"
                    }
                },
                {
                    "node": {
                        "name": "Modernize automation infrastructure and expand team Go competencies",
                        "status": {
                            "value": "pending"
                        },
                        "targetDate": {
                            "label": "August"
                        },
                        "id": "ari:cloud:townsquare:ef524de6-fdee-4937-a28b-5d531545cbae:goal/a208421e-a713-43e4-b110-904417ff3cc4"
                    }
                }
            ],
            "pageInfo": {
                "hasNextPage": false,
                "endCursor": "dHFsY29ubmVjdGlvbjoyOjUwOjA="
            }
        }
    }
}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, fakeResponse)
	}))
	defer server.Close()

	client := atlassian.NewClient("fake-cloud-id", "fake@email.com", "fake-token", server.URL)
	actualResult, err := client.FetchGoals(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(actualResult) != 3 {
		t.Errorf("expected 3 goals, got %d", len(actualResult))
	}
	if actualResult[0].Name != "First Go service deployed to GCP" {
		t.Errorf("unexpected name: %s", actualResult[0].Name)
	}
}

func TestFetchGoalsWithPagination(t *testing.T) {
	fakeResponsePage1 := `{
    "data": {
        "goals_search": {
            "edges": [
                {
                    "node": {
                        "name": "First Go service deployed to GCP",
                        "status": {
                            "value": "pending"
                        },
                        "targetDate": {
                            "label": "August"
                        },
                        "id": "ari:cloud:townsquare:ef524de6-fdee-4937-a28b-5d531545cbae:goal/a1a46457-943c-4df3-bf05-744ab089e271"
                    }
                }
            ],
            "pageInfo": {
                "hasNextPage": true,
                "endCursor": "dHFsY29ubmVjdGlvbjoyOjUwOjA="
            }
        }
    }
}`
	fakeResponsePage2 := `{
    "data": {
        "goals_search": {
            "edges": [
                {
                    "node": {
                        "name": "Second Goal",
                        "status": {
                            "value": "pending"
                        },
                        "targetDate": {
                            "label": "August"
                        },
                        "id": "ari:cloud:townsquare:ef524de6-fdee-4937-a28b-5d531545cbae:goal/2"
                    }
                }
            ],
            "pageInfo": {
                "hasNextPage": false,
                "endCursor": "abc"
            }
        }
    }
}`
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		if requestCount == 1 {
			fmt.Fprint(w, fakeResponsePage1)
		} else {
			fmt.Fprint(w, fakeResponsePage2)
		}

	}))
	defer server.Close()

	client := atlassian.NewClient("fake-cloud-id", "fake@email.com", "fake-token", server.URL)
	actualResult, err := client.FetchGoals(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(actualResult) != 2 {
		t.Errorf("expected 2 goals, got %d", len(actualResult))
	}
}

func TestFetchGoals_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized")
	}))
	defer server.Close()

	client := atlassian.NewClient("fake-cloud-id", "fake@email.com", "bad-token", server.URL)
	_, err := client.FetchGoals(t.Context())
	if err == nil {
		t.Fatal("expected error for non-200 response, got nil")
	}
}

func TestFetchGoals_GraphQLError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"errors":[{"message":"field 'goals_search' requires scope read:goal:townsquare"}]}`)
	}))
	defer server.Close()

	client := atlassian.NewClient("fake-cloud-id", "fake@email.com", "fake-token", server.URL)
	_, err := client.FetchGoals(t.Context())
	if err == nil {
		t.Fatal("expected error for GraphQL errors field, got nil")
	}
}

func TestFetchProjects(t *testing.T) {
	fakeResponse := `{
    "data": {
        "projects_search": {
            "edges": [
                {
                    "node": {
                        "id": "ari:cloud:townsquare:1:project/1",
                        "name": "Atlas Sync",
                        "archived": false,
                        "description": { "what": "Sync tool", "why": "Automation" },
                        "dueDate": { "label": "September" }
                    }
                },
                {
                    "node": {
                        "id": "ari:cloud:townsquare:1:project/2",
                        "name": "Monorepo Migration",
                        "archived": false,
                        "description": { "what": "Migrate services", "why": "Simplify" },
                        "dueDate": { "label": "October" }
                    }
                }
            ],
            "pageInfo": { "hasNextPage": false, "endCursor": "abc" }
        }
    }
}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, fakeResponse)
	}))
	defer server.Close()

	client := atlassian.NewClient("fake-cloud-id", "fake@email.com", "fake-token", server.URL)
	actualResult, err := client.FetchProjects(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(actualResult) != 2 {
		t.Errorf("expected 2 projects, got %d", len(actualResult))
	}
	if actualResult[0].Name != "Atlas Sync" {
		t.Errorf("unexpected name: %s", actualResult[0].Name)
	}
}
