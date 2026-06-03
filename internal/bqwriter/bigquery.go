package bqwriter

import (
	"bytes"
	"context"
	"encoding/json"

	"cloud.google.com/go/bigquery"
	"github.com/RomanUrakhov/atlas-sync/internal/models"
)

var _ Writer = (*BigQueryWriter)(nil)

type BigQueryWriter struct {
	client  *bigquery.Client
	dataset string
}

func NewBigQueryWriter(ctx context.Context, projectID, dataset string) (*BigQueryWriter, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &BigQueryWriter{client: client, dataset: dataset}, nil
}

type goalRow struct {
	ID         string `json:"id"          bigquery:"id"`
	Name       string `json:"name"        bigquery:"name"`
	Status     string `json:"status"      bigquery:"status"`
	TargetDate string `json:"target_date" bigquery:"target_date"`
}

type projectRow struct {
	ID              string `json:"id"               bigquery:"id"`
	Name            string `json:"name"             bigquery:"name"`
	DescriptionWhat string `json:"description_what" bigquery:"description_what"`
	DescriptionWhy  string `json:"description_why"  bigquery:"description_why"`
	DueDate         string `json:"due_date"         bigquery:"due_date"`
	Archived        bool   `json:"archived"         bigquery:"archived"`
}

func insertRows[T any](ctx context.Context, table *bigquery.Table, rows []T) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	for _, row := range rows {
		if err := enc.Encode(row); err != nil {
			return err
		}
	}

	src := bigquery.NewReaderSource(&buf)
	src.SourceFormat = bigquery.JSON

	loader := table.LoaderFrom(src)
	loader.WriteDisposition = bigquery.WriteAppend

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	return status.Err()
}

func (w *BigQueryWriter) InsertGoals(ctx context.Context, goals []models.Goal) error {
	rows := make([]goalRow, len(goals))
	for i, g := range goals {
		rows[i] = goalRow{
			ID:         g.ID,
			Name:       g.Name,
			Status:     g.Status.Value,
			TargetDate: g.TargetDate.Label,
		}
	}
	return insertRows(ctx, w.client.Dataset(w.dataset).Table("goals"), rows)
}

func (w *BigQueryWriter) InsertProjects(ctx context.Context, projects []models.Project) error {
	rows := make([]projectRow, len(projects))
	for i, p := range projects {
		rows[i] = projectRow{
			ID:              p.ID,
			Name:            p.Name,
			DescriptionWhat: p.Description.What,
			DescriptionWhy:  p.Description.Why,
			DueDate:         p.DueDate.Label,
			Archived:        p.Archived,
		}
	}
	return insertRows(ctx, w.client.Dataset(w.dataset).Table("projects"), rows)
}
