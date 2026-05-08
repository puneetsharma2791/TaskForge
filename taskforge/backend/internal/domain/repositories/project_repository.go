package repositories

import (
	"fmt"

	"github.com/acme/taskforge/internal/domain/aggregates"
	"github.com/acme/taskforge/internal/storage"
)

type ProjectRepository interface {
	Save(project *aggregates.Project) error
	FindByID(id string) (*aggregates.Project, error)
	Delete(id string) error
	ListByTenant(tenantID string) ([]*aggregates.Project, error)
}

type InMemoryProjectRepository struct {
	store     *storage.EventStore
	snapshots map[string]*aggregates.Project
}

func NewInMemoryProjectRepository(store *storage.EventStore) *InMemoryProjectRepository {
	return &InMemoryProjectRepository{
		store:     store,
		snapshots: make(map[string]*aggregates.Project),
	}
}

func (r *InMemoryProjectRepository) Save(project *aggregates.Project) error {
	pending := project.PendingEvents()
	if len(pending) == 0 {
		return nil
	}
	if err := r.store.Save(project.ID(), pending); err != nil {
		return err
	}
	r.snapshots[project.ID()] = project
	project.ClearEvents()
	return nil
}

func (r *InMemoryProjectRepository) FindByID(id string) (*aggregates.Project, error) {
	if p, ok := r.snapshots[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("project not found: %s", id)
}

func (r *InMemoryProjectRepository) Delete(id string) error {
	delete(r.snapshots, id)
	return nil
}

func (r *InMemoryProjectRepository) ListByTenant(tenantID string) ([]*aggregates.Project, error) {
	var result []*aggregates.Project
	for _, p := range r.snapshots {
		if p.TenantID() == tenantID {
			result = append(result, p)
		}
	}
	return result, nil
}
