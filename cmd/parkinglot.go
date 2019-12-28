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
func (pl *ParkingLot) park(registrationNumber string, color string) (*Slot, error) {
	if err := pl.isCreated(); err != nil {
		return nil, err
	}
	slotNumber, err := pl.getNearestParkingSlot()
	if err != nil {
		return nil, err
	}
	pl.slots[slotNumber-1].parkVehicle(createVehicle(registrationNumber, color))

	return pl.slots[slotNumber-1], nil
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

// Remove vehicle from parking slot
func (pl *ParkingLot) leave(slotNumber int) error {
	if err := pl.isCreated(); err != nil {
		return err
	}
	if slotNumber <= 0 || slotNumber > pl.capacity {
		return errors.New("Invalid slot number")
	}

	slot := pl.slots[slotNumber-1]
	if slot.getVehicle() != nil {
		// Remove vehicle from slot
		slot.removeVehicle()
		// Add empty slot to the heap
		heap.Push(&pl.emptySlot, &qheap.Item{Value: slotNumber})

		return nil
	}

	return errors.New("Vehicle is not found in parking lot")
}

// Get a list of vehicles parked in the parking lot, ordered by slot number
func (pl *ParkingLot) getStatus() []*Slot {
	if err := pl.isCreated(); err != nil {
		return nil
	}

	var slots []*Slot

	for i := 0; i < pl.highestSlot; i++ {
		slot := pl.slots[i]
		vehicle := slot.getVehicle()
		if vehicle != nil {
			slots = append(slots, slot)
		}
	}

	return slots
}

func (pl *ParkingLot) isCreated() error {
	if pl.capacity <= 0 {
		return errors.New("Parking lot is not created")
	}
	return nil
}
