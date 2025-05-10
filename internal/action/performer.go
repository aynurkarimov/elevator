package action

import (
	"fmt"

	"github.com/aynurkarimov/elevator/internal/elevator"
	"github.com/google/uuid"
)

type Command string

const (
	CommandReq Command = "req"
)

type Action struct {
	ElevatorId uuid.UUID
	Command    Command
	Floor      uint8
}

func PerformAction(action *Action) (int, error) {
	switch action.Command {
	case CommandReq:
		return elevator.RequestElevator(action.ElevatorId, action.Floor)
	default:
		return 0, fmt.Errorf("unknown action (%v)\n", action.Command)
	}
}
