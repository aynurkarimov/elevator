package elevator

import (
	"github.com/google/uuid"
)

var Elevators = []Elevator{
	{
		Id:           uuid.New(),
		CurrentFloor: 0,
	},
	{
		Id:           uuid.New(),
		CurrentFloor: 10,
	},
}
