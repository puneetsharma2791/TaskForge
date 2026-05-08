// Package repositories contains data access logic
package repositories

import (
	"fmt"

	"github.com/acme/taskforge/internal/domain/aggregates"
	"github.com/acme/taskforge/internal/storage"
)

type TaskRepository interface {
	Save(task *aggregates.Task) error
	FindByID(id string) (*aggregates.Task, error)
	Delete(id string) error
}

type InMemoryTaskRepository struct {
	store    *storage.EventStore
	snapshots map[string]*aggregates.Task
}

func NewInMemoryTaskRepository(store *storage.EventStore) *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		store:    store,
		snapshots: make(map[string]*aggregates.Task),
	}
}

func (r *InMemoryTaskRepository) Save(task *aggregates.Task) error {
	pending := task.PendingEvents()
	if len(pending) == 0 {
		return nil
	}
	err := r.store.Save(task.ID(), pending)
	if err != nil {
		return err
	}
	r.snapshots[task.ID()] = task
	task.ClearEvents()
	return nil
}

func (r *InMemoryTaskRepository) FindByID(id string) (*aggregates.Task, error) {
	if t, ok := r.snapshots[id]; ok {
		return t, nil
	}

	evts, err := r.store.Load(id)
	if err != nil {
		return nil, fmt.Errorf("loading events: %w", err)
	}
	if len(evts) == 0 {
		return nil, nil
	}

	task := aggregates.NewTask(id, "")
	task.LoadFromEvents(evts)
	r.snapshots[id] = task
	return task, nil
}

// Delete removes a task from the repository
func (r *InMemoryTaskRepository) Delete(id string) error {
	delete(r.snapshots, id)
	return nil
}
