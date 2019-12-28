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

// Park a vehicle
func (pl *ParkingLot) park(registrationNumber string, color string) (int, error) {
	if err := pl.isCreated(); err != nil {
		return 0, err
	}
	slotNumber, err := pl.getNearestParkingSlot()
	if err != nil {
		return 0, err
	}
	pl.slots[slotNumber-1].parkVehicle(createVehicle(registrationNumber, color))

	return slotNumber, nil
}

func (pl *ParkingLot) getNearestParkingSlot() (int, error) {
	var slotNumber int

	if pl.emptySlot.Len() == 0 {
		if pl.highestSlot == pl.capacity {
			return 0, errors.New("Sorry, parking lot is full")
		}
		slotNumber = pl.highestSlot + 1
		pl.highestSlot = slotNumber
	} else {
		item := heap.Pop(&pl.emptySlot)
		slotNumber = item.(*qheap.Item).Value
	}

	return slotNumber, nil
}

func (pl *ParkingLot) isCreated() error {
	if pl.capacity <= 0 {
		return errors.New("Parking lot is not created")
	}
	return nil
}
