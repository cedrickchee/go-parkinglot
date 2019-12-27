package cmd

import (
	"container/heap"
	"errors"

	qheap "github.com/cedrickchee/go-parkinglot/internal/heap"
)

type ParkingLot struct {
	address     string
	emptySlot   qheap.PriorityQueue
	slots       []*Slot
	highestSlot int
	capacity    int // Maximum slots available
}

type Slot struct {
	vehicle    *Vehicle
	slotNumber int
}

type Vehicle struct {
	registrationNumber string
	color              string
}

// Create parking lot
func (pl *ParkingLot) createParkingLot(address string, capacity int) error {
	if err := pl.isCreated(); err == nil {
		return errors.New("Parking lot already created")
	}
	pl.address = address
	pl.capacity = capacity

	var slots []*Slot
	for i := 0; i < capacity; i++ {
		slots = append(slots, &Slot{slotNumber: i + 1})
	}
	pl.slots = slots

	pl.emptySlot = qheap.PriorityQueue{}
	heap.Init(&pl.emptySlot) // Initialize the heap of empty slots

	return nil
}

func (pl *ParkingLot) isCreated() error {
	if pl.capacity <= 0 {
		return errors.New("Parking lot is not created")
	}
	return nil
}
