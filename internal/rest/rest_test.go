package rest

import (
	"net/http/httptest"
	"testing"

	"github.com/d-sauer/exploring-go/desk-controller/internal/domain/deskcontroller"
	mock_services "github.com/d-sauer/exploring-go/desk-controller/internal/domain/services/mock"
	"github.com/d-sauer/exploring-go/desk-controller/pkg/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_deskControllerRequestApi_GetControllers(t *testing.T) {
	mockService := mock_services.NewMockDeskControllerService(gomock.NewController(t))
	server := NewDeskControllerRequestAPI(mockService)
	t.Run("when successful", func(t *testing.T) {
		listOfControllers := [4]deskcontroller.Controller{
			{Identifier: "bup", Description: "Button up", State: deskcontroller.Idle},
			{Identifier: "bdown", Description: "Button down", State: deskcontroller.Idle},
			{Identifier: "bdone", Description: "Button one", State: deskcontroller.Idle},
			{Identifier: "bdtwo", Description: "Button two", State: deskcontroller.Idle},
		}

		mockService.EXPECT().List().Return(listOfControllers)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/v1/controllers", nil)

		server.GetControllers(rr, req)

		t.Run("it return 200", func(t *testing.T) {
			assert.Equal(t, 200, rr.Result().StatusCode)
		})
		t.Run("it matches OpenAPI", func(t *testing.T) {
			_, err := api.GetSwagger()
			assert.NoError(t, err)

			//_ = validator.NewValidator(doc).ForTest(t, rr, req)
		})
	})
}

func Test_deskControllerRequestApi_ControllerAction(t *testing.T) {
	// TODO
}

func Test_deskControllerRequestApi_GetControllerStatus(t *testing.T) {
	// TODO
}

func Test_deskControllerRequestApi_GetServiceHealth(t *testing.T) {
	// TODO
}

func TestDeskControllerRouter(t *testing.T) {
	// TODO
}
