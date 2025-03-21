package router

import (
	"net/http"

	"github.com/SkorikovGeorge/jobWorker/internal/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	return router
}

func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/jobs/pause", handlers.Pause).Methods(http.MethodGet).Name("Pause")
	router.HandleFunc("/jobs/resume", handlers.Resume).Methods(http.MethodGet).Name("Resume")
	router.HandleFunc("/jobs/{job_id}", handlers.GetJob).Methods(http.MethodGet).Name("GetJob")
	router.HandleFunc("/jobs", handlers.CreateJob).Methods(http.MethodPost).Name("CreateJob")
}
