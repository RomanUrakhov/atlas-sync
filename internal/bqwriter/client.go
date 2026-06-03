package bqwriter

import (
	"context"

	"github.com/RomanUrakhov/atlas-sync/internal/models"
)

type Writer interface {
	InsertGoals(ctx context.Context, goals []models.Goal) error
	InsertProjects(ctx context.Context, projects []models.Project) error
}
