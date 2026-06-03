package bqwriter

import (
	"context"

	"github.com/RomanUrakhov/atlas-sync/internal/models"
)

type MemoryWriter struct {
	Goals    []models.Goal
	Projects []models.Project
}

func (w *MemoryWriter) InsertGoals(_ context.Context, g []models.Goal) error {
	w.Goals = append(w.Goals, g...)
	return nil
}

func (w *MemoryWriter) InsertProjects(_ context.Context, p []models.Project) error {
	w.Projects = append(w.Projects, p...)
	return nil
}

func NewMemoryWriter() *MemoryWriter {
	return &MemoryWriter{}
}
