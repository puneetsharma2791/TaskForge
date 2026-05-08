package queries

import (
	"fmt"

	"github.com/acme/taskforge/internal/domain/projectors"
	"github.com/acme/taskforge/internal/domain/types"
)

type FindTask struct {
	TaskID   string
	TenantID string
}

func (q *FindTask) Execute(proj *projectors.TaskProjector) (*types.TaskView, error) {
	v := proj.GetView(q.TaskID)
	if v == nil {
		return nil, fmt.Errorf("task not found: %s", q.TaskID)
	}
	// TODO: tenant check
	return v, nil
}
