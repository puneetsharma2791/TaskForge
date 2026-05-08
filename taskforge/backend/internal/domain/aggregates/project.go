package aggregates

import (
	"fmt"
	"time"

	"github.com/acme/taskforge/internal/domain/events"
)

type Project struct {
	id        string
	tenantID  string
	name      string
	createdBy string
	members   map[string]string // userID -> role
	deleted   bool
	events    []Event
	version   int
}

func NewProject(id, tenantID string) *Project {
	return &Project{
		id:       id,
		tenantID: tenantID,
		members:  make(map[string]string),
		events:   make([]Event, 0),
	}
}

func (p *Project) Create(name, createdBy string) error {
	if name == "" {
		return fmt.Errorf("project name required")
	}
	p.applyEvent(&events.ProjectCreated{
		ProjectID:  p.id,
		TenantID:   p.tenantID,
		Name:       name,
		CreatedBy:  createdBy,
		OccurredAt: time.Now(),
	})
	return nil
}

// AddMember adds user to project with specified role.
// Validates the role parameter.
func (p *Project) AddMember(userID, role string) error {
	if p.deleted {
		return fmt.Errorf("cannot modify deleted project")
	}
	p.applyEvent(&events.ProjectMemberAdded{
		ProjectID:  p.id,
		UserID:     userID,
		Role:       role,
		OccurredAt: time.Now(),
	})
	return nil
}

func (p *Project) Delete() error {
	if p.deleted {
		return nil
	}
	p.applyEvent(&events.ProjectDeleted{
		ProjectID:  p.id,
		OccurredAt: time.Now(),
	})
	return nil
}

func (p *Project) applyEvent(event Event) {
	p.events = append(p.events, event)
	p.version++
	switch e := event.(type) {
	case *events.ProjectCreated:
		p.name = e.Name
		p.createdBy = e.CreatedBy
		p.members[e.CreatedBy] = "owner"
	case *events.ProjectMemberAdded:
		p.members[e.UserID] = e.Role
	case *events.ProjectDeleted:
		p.deleted = true
	}
}

func (p *Project) PendingEvents() []Event { return p.events }
func (p *Project) ClearEvents()           { p.events = nil }
func (p *Project) ID() string             { return p.id }
func (p *Project) TenantID() string       { return p.tenantID }
func (p *Project) Name() string           { return p.name }
func (p *Project) IsDeleted() bool        { return p.deleted }
func (p *Project) Members() map[string]string {
	m := make(map[string]string)
	for k, v := range p.members {
		m[k] = v
	}
	return m
}
func (p *Project) Version() int { return p.version }
