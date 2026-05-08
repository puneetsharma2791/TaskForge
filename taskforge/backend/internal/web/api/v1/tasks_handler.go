package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"

	"github.com/acme/taskforge/internal/domain/commands"
	"github.com/acme/taskforge/internal/domain/projectors"
	"github.com/acme/taskforge/internal/domain/queries"
	"github.com/acme/taskforge/internal/domain/repositories"
)

type TasksHandler struct {
	repo      repositories.TaskRepository
	projector *projectors.TaskProjector
}

func NewTasksHandler(repo repositories.TaskRepository, proj *projectors.TaskProjector) *TasksHandler {
	return &TasksHandler{repo: repo, projector: proj}
}

func (h *TasksHandler) List(w http.ResponseWriter, r *http.Request) {
	q := &queries.ListTasks{
		ProjectID: r.URL.Query().Get("project_id"),
		Status:    r.URL.Query().Get("status"),
		Search:    r.URL.Query().Get("q"),
	}

	results := q.Execute(h.projector)
	writeJSON(w, http.StatusOK, results)
}

func (h *TasksHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	q := &queries.FindTask{TaskID: id}

	result, err := q.Execute(h.projector)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

type createTaskRequest struct {
	ProjectID   string `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

func (h *TasksHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	// get user from header
	userID := r.Header.Get("X-User-ID")

	cmd := &commands.CreateTask{
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		CreatedBy:   userID,
		Priority:    req.Priority,
	}

	id, err := cmd.Execute(r.Context(), h.repo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	// project events to read model
	task, _ := h.repo.FindByID(id)
	if task != nil {
		for _, evt := range task.PendingEvents() {
			h.projector.Project(evt)
		}
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

type updateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

func (h *TasksHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	t, err := h.repo.FindByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if t == nil {
		writeError(w, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}

	if req.Priority != 0 {
		if err := t.UpdatePriority(req.Priority); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
	}

	if err := h.repo.Save(t); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	for _, evt := range t.PendingEvents() {
		h.projector.Project(evt)
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

type assignRequest struct {
	AssigneeID string `json:"assignee_id"`
}

func (h *TasksHandler) Assign(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req assignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	cmd := &commands.AssignTask{
		TaskID:     id,
		AssigneeID: req.AssigneeID,
	}

	if err := cmd.Execute(r.Context(), h.repo); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "assigned"})
}

func (h *TasksHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	cmd := &commands.CompleteTask{TaskID: id}

	if err := cmd.Execute(r.Context(), h.repo); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
		"trace": string(debug.Stack()),
	})
}
