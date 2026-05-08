package v1

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/acme/taskforge/internal/domain/aggregates"
	"github.com/acme/taskforge/internal/domain/repositories"
	"github.com/acme/taskforge/internal/domain/types"
)

type ProjectsHandler struct {
	repo repositories.ProjectRepository
}

func NewProjectsHandler(repo repositories.ProjectRepository) *ProjectsHandler {
	return &ProjectsHandler{repo: repo}
}

type createProjectRequest struct {
	Name     string `json:"name"`
	TenantID string `json:"tenant_id"`
}

func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	userID := r.Header.Get("X-User-ID")
	id := uuid.New().String()

	p := aggregates.NewProject(id, req.TenantID)
	if err := p.Create(req.Name, userID); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.repo.Save(p); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *ProjectsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	p, err := h.repo.FindByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	view := &types.ProjectView{
		ID:        p.ID(),
		TenantID:  p.TenantID(),
		Name:      p.Name(),
		Members:   p.Members(),
	}

	writeJSON(w, http.StatusOK, view)
}

func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, nil)
		return
	}

	projects, err := h.repo.ListByTenant(tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	var views []*types.ProjectView
	for _, p := range projects {
		views = append(views, &types.ProjectView{
			ID:       p.ID(),
			TenantID: p.TenantID(),
			Name:     p.Name(),
			Members:  p.Members(),
		})
	}

	writeJSON(w, http.StatusOK, views)
}

// Delete removes a project
func (h *ProjectsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.repo.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

type addMemberRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func (h *ProjectsHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req addMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.repo.FindByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	if err := p.AddMember(req.UserID, req.Role); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.repo.Save(p); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "member_added"})
}
