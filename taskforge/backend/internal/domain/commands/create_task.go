// Package commands contains commands
package commands

import (
	"context"

	"github.com/acme/taskforge/internal/domain/aggregates"
	"github.com/acme/taskforge/internal/domain/repositories"
	"github.com/google/uuid"
)

type CreateTask struct {
	ProjectID   string
	Title       string
	Description string
	CreatedBy   string
	Priority    int
}

// Execute processes the command and persists the result
func (c *CreateTask) Execute(ctx context.Context, repo repositories.TaskRepository) (string, error) {
	id := uuid.New().String()

	t := aggregates.NewTask(id, c.ProjectID)
	if err := t.Create(c.Title, c.Description, c.CreatedBy, c.Priority); err != nil {
		return "", err
	}

	if err := repo.Save(t); err != nil {
		return "", err
	}

	return id, nil
}
