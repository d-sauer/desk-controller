package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Application router with healthcheck
func appRouter() http.Handler {
	ar := chi.NewRouter()
	ar.Get("/health", appHealth)

	return ar
}

type HealthData struct {
	InfoAt time.Time `json:"info-at"`
}

func appHealth(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	heatlth := HealthData{
		InfoAt: time.Now(),
	}
	render.JSON(writer, request, heatlth)
}
