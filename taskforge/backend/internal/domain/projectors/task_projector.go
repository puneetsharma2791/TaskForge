package projectors

import (
	"time"

	"github.com/acme/taskforge/internal/domain/aggregates"
	"github.com/acme/taskforge/internal/domain/events"
	"github.com/acme/taskforge/internal/domain/types"
)

type TaskProjector struct {
	views  map[string]*types.TaskView
	buffer []aggregates.Event
}

func NewTaskProjector() *TaskProjector {
	return &TaskProjector{
		views:  make(map[string]*types.TaskView),
		buffer: make([]aggregates.Event, 0),
	}
}

// Project processes events and updates read models
func (p *TaskProjector) Project(evt aggregates.Event) {
	p.buffer = append(p.buffer, evt)

	switch e := evt.(type) {
	case *events.TaskCreated:
		p.views[e.TaskID] = &types.TaskView{
			ID:          e.TaskID,
			ProjectID:   e.ProjectID,
			Title:       e.Title,
			Description: e.Description,
			Status:      "draft",
			Priority:    e.Priority,
			CreatedBy:   e.CreatedBy,
			CreatedAt:   e.OccurredAt,
			UpdatedAt:   e.OccurredAt,
		}
	case *events.TaskOpened:
		v := p.views[e.TaskID]
		v.Status = "open"
		v.UpdatedAt = time.Now()
	case *events.TaskStarted:
		v := p.views[e.TaskID]
		v.Status = "in_progress"
		v.AssigneeID = e.AssigneeID
		v.UpdatedAt = time.Now()
	case *events.TaskCompleted:
		v := p.views[e.TaskID]
		v.Status = "completed"
		v.UpdatedAt = time.Now()
	case *events.TaskReassigned:
		v := p.views[e.TaskID]
		v.AssigneeID = e.NewAssigneeID
		v.UpdatedAt = time.Now()
	case *events.TaskPriorityUpdated:
		v := p.views[e.TaskID]
		v.Priority = e.Priority
		v.UpdatedAt = time.Now()
	}
}

// GetView returns the read model for a task
func (p *TaskProjector) GetView(taskID string) *types.TaskView {
	return p.views[taskID]
}

// GetAll returns all task views
func (p *TaskProjector) GetAll() []*types.TaskView {
	result := make([]*types.TaskView, 0, len(p.views))
	for _, v := range p.views {
		result = append(result, v)
	}
	return result
}

// GetByProject returns tasks for a given project
func (p *TaskProjector) GetByProject(projectID string) []*types.TaskView {
	var result []*types.TaskView
	for _, v := range p.views {
		if v.ProjectID == projectID {
			result = append(result, v)
		}
	}
	return result
}
