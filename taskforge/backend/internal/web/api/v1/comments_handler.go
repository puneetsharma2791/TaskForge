package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/acme/taskforge/internal/domain/projectors"
	"github.com/acme/taskforge/internal/domain/repositories"
	"github.com/acme/taskforge/internal/storage"
)

type CommentsHandler struct {
	repo      repositories.TaskRepository
	projector *projectors.TaskProjector
	store     *storage.EventStore
}

func NewCommentsHandler(repo repositories.TaskRepository, proj *projectors.TaskProjector, store *storage.EventStore) *CommentsHandler {
	return &CommentsHandler{repo: repo, projector: proj, store: store}
}

func (h *CommentsHandler) projectEventsForAggregate(aggregateID string) {
	events, _ := h.store.Load(aggregateID)
	for _, evt := range events {
		h.projector.Project(evt)
	}
}

// List returns all comments for a task
func (h *CommentsHandler) List(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]

	// Verify task exists
	view := h.projector.GetView(taskID)
	if view == nil {
		writeError(w, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}

	comments := h.projector.GetComments(taskID)
	if comments == nil {
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

type addCommentRequest struct {
	Content string `json:"content"`
}

// Add creates a new comment on a task
func (h *CommentsHandler) Add(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]
	userID := r.Header.Get("X-User-ID")

	var req addCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	t, err := h.repo.FindByID(taskID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if t == nil {
		writeError(w, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}

	commentID, err := t.AddComment(userID, req.Content)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	pendingEvents := t.PendingEvents()
	if err := h.repo.Save(t); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	for _, evt := range pendingEvents {
		h.projector.Project(evt)
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": commentID})
}

type editCommentRequest struct {
	Content string `json:"content"`
}

// Edit updates a comment's content
func (h *CommentsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]
	commentID := mux.Vars(r)["commentId"]
	userID := r.Header.Get("X-User-ID")

	var req editCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	t, err := h.repo.FindByID(taskID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if t == nil {
		writeError(w, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}

	if err := t.EditComment(commentID, req.Content, userID); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	pendingEvents := t.PendingEvents()
	if err := h.repo.Save(t); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	for _, evt := range pendingEvents {
		h.projector.Project(evt)
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// Delete removes a comment from a task
func (h *CommentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]
	commentID := mux.Vars(r)["commentId"]
	userID := r.Header.Get("X-User-ID")

	t, err := h.repo.FindByID(taskID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if t == nil {
		writeError(w, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}

	if err := t.DeleteComment(commentID, userID); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	pendingEvents := t.PendingEvents()
	if err := h.repo.Save(t); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	for _, evt := range pendingEvents {
		h.projector.Project(evt)
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
