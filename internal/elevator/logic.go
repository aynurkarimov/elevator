package elevator

import (
	"fmt"
	"math"
	"time"

	"github.com/aynurkarimov/elevator/internal/constants"
	"github.com/google/uuid"
)

func findElevator(elevatorId uuid.UUID) (*Elevator, error) {
	for i := range Elevators {
		if Elevators[i].Id == elevatorId {
			return &Elevators[i], nil
		}
	}

	return nil, fmt.Errorf("couldn't find %v elevator\n", elevatorId)
}

// FIXME: channel seems out of place
func runElevator(teardownCh chan int, elevator *Elevator, requestedFloor uint8) {
	if requestedFloor == elevator.CurrentFloor {
		return
	}

	elevator.Queue = append(elevator.Queue, requestedFloor)

	floors := int(math.Abs(float64(int(elevator.CurrentFloor) - int(requestedFloor))))
	difference := int((constants.SPEED * time.Duration(floors)).Milliseconds())

	teardownCh <- difference

	if requestedFloor > elevator.CurrentFloor {
		for i := elevator.CurrentFloor; i < requestedFloor; i++ {
			time.Sleep(constants.SPEED)

			elevator.Mu.Lock()
			elevator.CurrentFloor += 1
			elevator.Mu.Unlock()
		}
	} else {
		for i := elevator.CurrentFloor; i > requestedFloor; i-- {
			time.Sleep(constants.SPEED)

			elevator.Mu.Lock()
			elevator.CurrentFloor -= 1
			elevator.Mu.Unlock()
		}
	}

	elevator.Queue = elevator.Queue[1:]
}

func RequestElevator(elevatorId uuid.UUID, requestedFloor uint8) (int, error) {
	elevator, err := findElevator(elevatorId)

	if err != nil {
		return 0, err
	}

	teardownCh := make(chan int)

	go runElevator(teardownCh, elevator, requestedFloor)

	ms := <-teardownCh

	return ms, nil
}
