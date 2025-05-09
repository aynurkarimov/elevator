package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Command string

const (
	CommandReq Command = "req"
)

type Action struct {
	elevatorId uuid.UUID
	command    Command
	floor      uint8
}

type Elevator struct {
	id           uuid.UUID
	currentFloor uint8
	mu           sync.Mutex
}

var elevators = make([]Elevator, 1)

const SPEED = time.Second

func parseUserAction() (*Action, error) {
	reader := bufio.NewReader(os.Stdin)

	str, err := reader.ReadString('\n')

	if err != nil {
		return nil, fmt.Errorf("couldn't read user's input\n")
	}

	parts := strings.Split(str, ",")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid format for action (less than two)\n")
	}

	// FIXME: must be dynamic
	if strings.TrimSpace(parts[0]) != string(CommandReq) {
		return nil, fmt.Errorf("invalid format for action (wrong command)\n")
	}

	floor, err := strconv.Atoi(strings.TrimSpace(parts[1]))

	if err != nil {
		return nil, nil
	}

	if floor < 0 || floor > 255 {
		return nil, nil
	}

	return &Action{
		// elevatorId: uuid.UUID, // FIXME
		command: CommandReq,
		floor:   uint8(floor),
	}, nil
}

func tick(elevator *Elevator) {
	fmt.Printf("id: %v\n", elevator.id)
	fmt.Printf("floor: %v\n", elevator.currentFloor)
}

func findElevator(elevatorId uuid.UUID) (*Elevator, error) {
	for i := range elevators {
		if elevators[i].id == elevatorId {
			return &elevators[i], nil
		}
	}

	return nil, fmt.Errorf("couldn't find %v elevator\n", elevatorId)
}

func runElevator(elevator *Elevator, requestedFloor uint8) {
	tick(elevator)

	if requestedFloor == elevator.currentFloor {
		return
	}

	if requestedFloor > elevator.currentFloor {
		for i := elevator.currentFloor; i < requestedFloor; i++ {
			time.Sleep(SPEED)

			elevator.mu.Lock()

			elevator.currentFloor += 1

			tick(elevator)

			elevator.mu.Unlock()
		}

	} else {
		for i := elevator.currentFloor; i > requestedFloor; i-- {
			//
		}
	}
}

func requestElevator(elevatorId uuid.UUID, requestedFloor uint8) error {
	elevator, err := findElevator(elevatorId)

	if err != nil {
		return err
	}

	go runElevator(elevator, requestedFloor)

	return nil
}

func performUserAction(action *Action) error {
	switch action.command {
	case CommandReq:
		return requestElevator(action.elevatorId, action.floor)
	default:
		return fmt.Errorf("unknown action (%v)\n", action.command)
	}
}

func main() {
	for {
		action, err := parseUserAction()

		if err != nil {
			panic(err)
		}

		err = performUserAction(action)

		if err != nil {
			panic(err)
		}
	}
}
