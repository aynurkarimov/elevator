package elevator

import (
	"sync"

	"github.com/google/uuid"
)

type Elevator struct {
	Id           uuid.UUID
	CurrentFloor uint8
	Queue        []uint8
	Mu           sync.Mutex
}
