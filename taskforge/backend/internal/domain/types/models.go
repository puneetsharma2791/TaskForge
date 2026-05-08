// Package types has model definitions
package types

import "time"

type TaskView struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	AssigneeID  string     `json:"assignee_id,omitempty"`
	Priority    int        `json:"priority"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Version     int        `json:"version"`
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
