package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/acme/taskforge/internal/config"
	"github.com/acme/taskforge/internal/domain/projectors"
	"github.com/acme/taskforge/internal/domain/repositories"
	"github.com/acme/taskforge/internal/storage"
	v1 "github.com/acme/taskforge/internal/web/api/v1"
)

func main() {
	cfg := config.Load()

	store := storage.NewEventStore()
	taskRepo := repositories.NewInMemoryTaskRepository(store)
	projectRepo := repositories.NewInMemoryProjectRepository(store)
	taskProjector := projectors.NewTaskProjector()

	tasksHandler := v1.NewTasksHandler(taskRepo, taskProjector, store)
	projectsHandler := v1.NewProjectsHandler(projectRepo)
	commentsHandler := v1.NewCommentsHandler(taskRepo, taskProjector, store)

	r := mux.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-Tenant-ID")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	v1.RegisterRoutes(r, tasksHandler, projectsHandler, commentsHandler)

	// Wrap the router with a top-level CORS handler so OPTIONS
	// preflight requests are caught before mux route matching.
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-Tenant-ID")
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		r.ServeHTTP(w, req)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("TaskForge API starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
