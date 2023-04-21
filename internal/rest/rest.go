package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/d-sauer/exploring-go/desk-controller/internal/domain/services"
	"github.com/d-sauer/exploring-go/desk-controller/pkg/api"
	"github.com/go-chi/chi/v5"
)

type deskControllerRequestApi struct {
	service services.DeskControllerService
}

func NewDeskControllerRequestAPI(service services.DeskControllerService) api.ServerInterface {
	return &deskControllerRequestApi{
		service: service,
	}
}

func (ra *deskControllerRequestApi) GetServiceHealth(w http.ResponseWriter, r *http.Request) {

	// TODO: Implement
}

func (ra *deskControllerRequestApi) GetControllers(w http.ResponseWriter, r *http.Request) {

	// TODO: Implement
}

func (ra *deskControllerRequestApi) GetControllerStatus(w http.ResponseWriter, r *http.Request, controllerSlug string) {

	// TODO: Implement
}

func (ra *deskControllerRequestApi) ControllerAction(w http.ResponseWriter, r *http.Request, controllerSlug string) {

	// TODO: Implement
}

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
