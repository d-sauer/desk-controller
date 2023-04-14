package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/d-sauer/exploring-go/desk-controller/pkg/api"
	"github.com/go-chi/chi/v5"
)

func DeskControllerRouter() http.Handler {
	openapi, err := api.GetSwagger()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	openapi.Servers = nil

	r := chi.NewRouter()

	return r
}
