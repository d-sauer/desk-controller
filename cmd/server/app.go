package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Application router with healthcheck
func appRouter() http.Handler {
	ar := chi.NewRouter()
	ar.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(fmt.Sprintf("time: %s", time.Now())))
	})

	return ar
}
