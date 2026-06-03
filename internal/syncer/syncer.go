package syncer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/RomanUrakhov/atlas-sync/internal/bqwriter"
	"github.com/RomanUrakhov/atlas-sync/internal/models"
)

type Fetcher interface {
	FetchGoals(ctx context.Context) ([]models.Goal, error)
	FetchProjects(ctx context.Context) ([]models.Project, error)
}

type Syncer struct {
	fetcher Fetcher
	writer  bqwriter.Writer
}

func NewSyncer(f Fetcher, w bqwriter.Writer) *Syncer {
	return &Syncer{f, w}
}

// TODO: rewrite with concurrency later
func (s *Syncer) Run(ctx context.Context) error {
	slog.Info("fetching goals")
	goals, err := s.fetcher.FetchGoals(ctx)
	if err != nil {
		return fmt.Errorf("fetch goals: %w", err)
	}
	slog.Info("fetched goals", "count", len(goals))

	slog.Info("fetching projects")
	projects, err := s.fetcher.FetchProjects(ctx)
	if err != nil {
		return fmt.Errorf("fetch projects: %w", err)
	}
	slog.Info("fetched projects", "count", len(projects))

	slog.Info("uploading goals")
	if err := s.writer.InsertGoals(ctx, goals); err != nil {
		return fmt.Errorf("insert goals: %w", err)
	}
	slog.Info("uploaded goals")

	slog.Info("uploading projects")
	if err := s.writer.InsertProjects(ctx, projects); err != nil {
		return fmt.Errorf("insert projects: %w", err)
	}
	slog.Info("uploaded projects")

	return nil
}
