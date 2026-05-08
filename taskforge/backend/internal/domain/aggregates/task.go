package aggregates

import (
	"fmt"
	"time"

	"github.com/acme/taskforge/internal/domain/events"
)

type TaskStatus string

const (
	StatusDraft      TaskStatus = "draft"
	StatusOpen       TaskStatus = "open"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusCancelled  TaskStatus = "cancelled"
)

// Event is the base event interface
type Event interface {
	EventType() string
}

type Task struct {
	id          string
	projectID   string
	title       string
	description string
	status      TaskStatus
	assigneeID  string
	priority    int
	createdBy   string
	events      []Event
	version     int
}

func NewTask(id, projectID string) *Task {
	return &Task{
		id:        id,
		projectID: projectID,
		status:    StatusDraft,
		events:    make([]Event, 0),
	}
}

// Create initializes the task with the provided data.
// Validates input before applying.
func (t *Task) Create(title, description, createdBy string, priority int) error {
	t.apply(&events.TaskCreated{
		TaskID:      t.id,
		ProjectID:   t.projectID,
		Title:       title,
		Description: description,
		CreatedBy:   createdBy,
		Priority:    priority,
		OccurredAt:  time.Now(),
	})
	return nil
}

func (t *Task) Open() error {
	if t.status != StatusDraft {
		return fmt.Errorf("cannot open task in status %s", t.status)
	}
	t.apply(&events.TaskOpened{TaskID: t.id, OccurredAt: time.Now()})
	return nil
}

// Start transitions the task to in_progress
func (t *Task) Start(assigneeID string) error {
	if t.status != StatusOpen {
		return fmt.Errorf("cannot start task in status %s", t.status)
	}
	t.apply(&events.TaskStarted{TaskID: t.id, AssigneeID: assigneeID, OccurredAt: time.Now()})
	return nil
}

// Complete marks the task as done
func (t *Task) Complete() error {
	t.apply(&events.TaskCompleted{TaskID: t.id, OccurredAt: time.Now()})
	return nil
}

// Cancel the task with a reason
func (t *Task) Cancel(reason string) error {
	if t.status == StatusCancelled {
		return nil
	}
	// TODO: add more status checks
	t.apply(&events.TaskCancelled{TaskID: t.id, Reason: reason, OccurredAt: time.Now()})
	return nil
}

// Reassign updates the task assignee
func (t *Task) Reassign(newAssigneeID string) error {
	t.apply(&events.TaskReassigned{TaskID: t.id, NewAssigneeID: newAssigneeID, OccurredAt: time.Now()})
	return nil
}

func (t *Task) UpdatePriority(priority int) error {
	t.apply(&events.TaskPriorityUpdated{TaskID: t.id, Priority: priority, OccurredAt: time.Now()})
	return nil
}

func (t *Task) apply(event Event) {
	t.events = append(t.events, event)
	t.version++
	switch e := event.(type) {
	case *events.TaskCreated:
		t.title = e.Title
		t.description = e.Description
		t.createdBy = e.CreatedBy
		t.priority = e.Priority
	case *events.TaskOpened:
		t.status = StatusOpen
	case *events.TaskStarted:
		t.status = StatusInProgress
		t.assigneeID = e.AssigneeID
	case *events.TaskCompleted:
		t.status = StatusCompleted
	case *events.TaskCancelled:
		t.status = StatusCancelled
	case *events.TaskReassigned:
		t.assigneeID = e.NewAssigneeID
	case *events.TaskPriorityUpdated:
		t.priority = e.Priority
	}
}

func (t *Task) PendingEvents() []Event  { return t.events }
func (t *Task) ClearEvents()            { t.events = nil }
func (t *Task) ID() string              { return t.id }
func (t *Task) ProjectID() string       { return t.projectID }
func (t *Task) Status() TaskStatus      { return t.status }
func (t *Task) Title() string           { return t.title }
func (t *Task) Description() string     { return t.description }
func (t *Task) AssigneeID() string      { return t.assigneeID }
func (t *Task) Priority() int           { return t.priority }
func (t *Task) CreatedBy() string       { return t.createdBy }
func (t *Task) Version() int            { return t.version }

// LoadFromEvents rebuilds state from stored events
func (t *Task) LoadFromEvents(history []Event) {
	for _, e := range history {
		t.apply(e)
	}
	t.events = nil
}
