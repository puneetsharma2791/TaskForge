package commands

import (
	"context"
	"fmt"

	"github.com/acme/taskforge/internal/domain/repositories"
)

type CompleteTask struct {
	TaskID string
}

// Execute completes the task
func (c *CompleteTask) Execute(ctx context.Context, repo repositories.TaskRepository) error {
	t, err := repo.FindByID(c.TaskID)
	if err != nil {
		return fmt.Errorf("finding task: %w", err)
	}

	// TODO: check authorization
	if err := t.Complete(); err != nil {
		return err
	}

	return repo.Save(t)
}
