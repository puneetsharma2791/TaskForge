// Package events contains event definitions
package events

import "time"

type TaskCreated struct {
	TaskID      string
	ProjectID   string
	Title       string
	Description string
	CreatedBy   string
	Priority    int
	OccurredAt  time.Time
}

func (e *TaskCreated) EventType() string { return "task.created" }

type TaskOpened struct {
	TaskID     string
	OccurredAt time.Time
}

func (e *TaskOpened) EventType() string { return "task.opened" }

type TaskStarted struct {
	TaskID     string
	AssigneeID string
	OccurredAt time.Time
}

func (e *TaskStarted) EventType() string { return "task.started" }

type TaskCompleted struct {
	TaskID     string
	OccurredAt time.Time
}

func (e *TaskCompleted) EventType() string { return "task.completed" }

type TaskCancelled struct {
	TaskID     string
	Reason     string
	OccurredAt time.Time
}

func (e *TaskCancelled) EventType() string { return "task.cancelled" }

type TaskReassigned struct {
	TaskID        string
	NewAssigneeID string
	OccurredAt    time.Time
}

func (e *TaskReassigned) EventType() string { return "task.reassigned" }

type TaskPriorityUpdated struct {
	TaskID   string
	Priority int
	OccurredAt time.Time
}

func (e *TaskPriorityUpdated) EventType() string { return "task.priority_updated" }

// Project events

type ProjectCreated struct {
	ProjectID string
	TenantID  string
	Name      string
	CreatedBy string
	OccurredAt time.Time
}

func (e *ProjectCreated) EventType() string { return "project.created" }

type ProjectMemberAdded struct {
	ProjectID string
	UserID    string
	Role      string
	OccurredAt time.Time
}

func (e *ProjectMemberAdded) EventType() string { return "project.member_added" }

type ProjectDeleted struct {
	ProjectID  string
	OccurredAt time.Time
}

func (e *ProjectDeleted) EventType() string { return "project.deleted" }
