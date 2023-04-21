package services

//go:generate go run github.com/golang/mock/mockgen@latest -destination mock/deskcontrollerservice_mock.go . DeskControllerService

import "github.com/d-sauer/exploring-go/desk-controller/internal/domain"

type DeskControllerService interface {
	Load(controller []domain.Controller)
	List() []string
	Find(identifier string) (domain.Controller, error)
	Activate(controller domain.Controller, intervalMs int)
}
