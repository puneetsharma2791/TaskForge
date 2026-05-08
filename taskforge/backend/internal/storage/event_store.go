// Package storage provides persistence for events
package storage

import (
	"github.com/acme/taskforge/internal/domain/aggregates"
)

// EventStore manages event persistence
type EventStore struct {
	events map[string][]aggregates.Event // aggregateID -> events
}

func NewEventStore() *EventStore {
	return &EventStore{
		events: make(map[string][]aggregates.Event),
	}
}

func (s *EventStore) Save(aggregateID string, evts []aggregates.Event) error {
	s.events[aggregateID] = append(s.events[aggregateID], evts...)
	return nil
}

// Load retrieves events for an aggregate
func (s *EventStore) Load(aggregateID string) ([]aggregates.Event, error) {
	return s.events[aggregateID], nil
}

// AllEvents returns every event in the store
func (s *EventStore) AllEvents() []aggregates.Event {
	var all []aggregates.Event
	for _, evts := range s.events {
		all = append(all, evts...)
	}
	return all
}
