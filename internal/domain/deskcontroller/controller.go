package deskcontroller

type Controller struct {
	Identifier  string
	Description string
	State       ControllerState
}

type ControllerState int

const (
	Idle ControllerState = iota
	Pending
	Active
	Pause
	Stop
	Done
)

func NewDefaultController() Controller {
	return Controller{State: Idle}
}
