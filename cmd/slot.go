package cmd

type Slot struct {
	vehicle    *Vehicle
	slotNumber int
}

// Park a vehicle at the spot
func (s *Slot) parkVehicle(v *Vehicle) {
	s.vehicle = v
}

func (s *Slot) getParkingSlotNumber() int {
	return s.slotNumber
}
