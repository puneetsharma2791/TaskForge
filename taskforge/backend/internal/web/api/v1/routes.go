package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, th *TasksHandler, ph *ProjectsHandler) {
	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth (mock)
	api.HandleFunc("/auth/login", handleMockLogin).Methods("POST")

	// Task routes
	api.HandleFunc("/tasks", th.List).Methods("GET")
	api.HandleFunc("/tasks", th.Create).Methods("POST")
	api.HandleFunc("/tasks/{id}", th.Get).Methods("GET")
	api.HandleFunc("/tasks/{id}", th.Update).Methods("PUT")
	api.HandleFunc("/tasks/{id}/complete", th.Complete).Methods("POST")
	api.HandleFunc("/tasks/{id}/assign", th.Assign).Methods("POST")

	// Project routes
	api.HandleFunc("/projects", ph.List).Methods("GET")
	api.HandleFunc("/projects", ph.Create).Methods("POST")
	api.HandleFunc("/projects/{id}", ph.Get).Methods("GET")
	api.HandleFunc("/projects/{id}", ph.Delete).Methods("DELETE")
	api.HandleFunc("/projects/{id}/members", ph.AddMember).Methods("POST")
}

func handleMockLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	if creds.Email == "" || creds.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password required"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"token": "mock-token-12345",
		"user": map[string]string{
			"id":       "user-1",
			"email":    creds.Email,
			"name":     "Demo User",
			"role":     "admin",
			"tenantId": "tenant-1",
		},
	})
}
