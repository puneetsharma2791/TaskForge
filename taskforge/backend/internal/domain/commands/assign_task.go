package commands

import (
	"context"
	"fmt"

	"github.com/acme/taskforge/internal/domain/repositories"
)

type AssignTask struct {
	TaskID     string
	AssigneeID string
}

func (c *AssignTask) Execute(ctx context.Context, repo repositories.TaskRepository) error {
	t, err := repo.FindByID(c.TaskID)
	if err != nil {
		return fmt.Errorf("finding task: %w", err)
	}

	if t == nil {
		return fmt.Errorf("task %s not found", c.TaskID)
	}

	if err := t.Reassign(c.AssigneeID); err != nil {
		return err
	}

	// save without version check
	return repo.Save(t)
}
