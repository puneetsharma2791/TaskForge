// Package types has model definitions
package types

import "time"

type TaskView struct {
	ID           string     `json:"id"`
	ProjectID    string     `json:"project_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	AssigneeID   string     `json:"assignee_id,omitempty"`
	Priority     int        `json:"priority"`
	CreatedBy    string     `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Version      int        `json:"version"`
	CommentCount int        `json:"comment_count"`
}

type CommentView struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectView struct {
	ID        string            `json:"id"`
	TenantID  string            `json:"tenant_id"`
	Name      string            `json:"name"`
	CreatedBy string            `json:"created_by"`
	Members   map[string]string `json:"members"`
	CreatedAt time.Time         `json:"created_at"`
}

// Member roles
const (
	RoleOwner  = "owner"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)
