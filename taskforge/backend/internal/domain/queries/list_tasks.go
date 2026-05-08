package queries

import (
	"strings"

	"github.com/acme/taskforge/internal/domain/projectors"
	"github.com/acme/taskforge/internal/domain/types"
)

type ListTasks struct {
	ProjectID string
	Status    string
	Search    string
}

// Execute returns matching tasks from the projector.
// Supports filtering by project, status and free-text search.
func (q *ListTasks) Execute(proj *projectors.TaskProjector) []*types.TaskView {
	var results []*types.TaskView

	all := proj.GetAll()
	for _, t := range all {
		if q.ProjectID != "" && t.ProjectID != q.ProjectID {
			continue
		}
		if q.Status != "" && t.Status != q.Status {
			continue
		}
		if q.Search != "" && !matchSearch(t, q.Search) {
			continue
		}
		results = append(results, t)
	}

	return results
}

func matchSearch(t *types.TaskView, search string) bool {
	s := strings.ToLower(search)
	return strings.Contains(strings.ToLower(t.Title), s) ||
		strings.Contains(strings.ToLower(t.Description), s)
}
