package syncer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/RomanUrakhov/atlas-sync/internal/bqwriter"
	"github.com/RomanUrakhov/atlas-sync/internal/models"
	"github.com/RomanUrakhov/atlas-sync/internal/syncer"
)

type fakeFetcher struct {
	goals       []models.Goal
	projects    []models.Project
	goalsErr    error
	projectsErr error
}

func (f *fakeFetcher) FetchGoals(_ context.Context) ([]models.Goal, error) {
	return f.goals, f.goalsErr
}

func (f *fakeFetcher) FetchProjects(_ context.Context) ([]models.Project, error) {
	return f.projects, f.projectsErr
}

func TestSyncer_Run(t *testing.T) {
	fetcher := &fakeFetcher{
		goals:    []models.Goal{{ID: "g1", Name: "Goal One"}},
		projects: []models.Project{{ID: "p1", Name: "Project One"}},
	}
	writer := bqwriter.NewMemoryWriter()
	s := syncer.NewSyncer(fetcher, writer)

	if err := s.Run(t.Context()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(writer.Goals) != 1 || writer.Goals[0].ID != "g1" {
		t.Errorf("unexpected goals in writer: %v", writer.Goals)
	}
	if len(writer.Projects) != 1 || writer.Projects[0].ID != "p1" {
		t.Errorf("unexpected projects in writer: %v", writer.Projects)
	}
}

func TestSyncer_Run_FetchGoalsError(t *testing.T) {
	apiErr := errors.New("api down")
	s := syncer.NewSyncer(&fakeFetcher{goalsErr: apiErr}, bqwriter.NewMemoryWriter())

	err := s.Run(t.Context())
	if !errors.Is(err, apiErr) {
		t.Fatalf("expected wrapped apiErr, got: %v", err)
	}
}

func TestSyncer_Run_FetchProjectsError(t *testing.T) {
	apiErr := errors.New("api down")
	s := syncer.NewSyncer(
		&fakeFetcher{
			goals:       []models.Goal{{ID: "g1"}},
			projectsErr: apiErr,
		},
		bqwriter.NewMemoryWriter(),
	)

	err := s.Run(t.Context())
	if !errors.Is(err, apiErr) {
		t.Fatalf("expected wrapped apiErr, got: %v", err)
	}
}
