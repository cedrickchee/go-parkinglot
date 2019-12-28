package cmd

import (
	"reflect"
	"testing"

	qheap "github.com/cedrickchee/go-parkinglot/internal/heap"
)

func generateParkingSlot(capacity int) []*Slot {
	var slots []*Slot

	for i := 0; i < capacity; i++ {
		slots = append(slots, &Slot{slotNumber: i + 1})
	}
	return slots
}

func compareParkingLot(t *testing.T, got *ParkingLot, want *ParkingLot) {
	if !reflect.DeepEqual(got.emptySlot, want.emptySlot) ||
		!reflect.DeepEqual(got.slots, want.slots) ||
		got.address != want.address ||
		got.highestSlot != want.highestSlot ||
		got.capacity != want.capacity {
		t.Errorf("createParkingLot() got = %v, want = %v", got, want)
	}
}

func TestCreateParkingLot(t *testing.T) {
	slots := generateParkingSlot(10)
	address := "Marina Bay Sands"

	type args struct {
		address  string
		capacity int
	}
	tests := []struct {
		name       string
		parkinglot *ParkingLot
		args       args
		want       *ParkingLot
		wantErr    bool
	}{
		{
			name:       "Parking lot is not created",
			parkinglot: &ParkingLot{},
			args: args{
				address:  address,
				capacity: 10,
			},
			want:    &ParkingLot{address: address, emptySlot: qheap.PriorityQueue{}, slots: slots, highestSlot: 0, capacity: 10},
			wantErr: false,
		},
		{
			name:       "Parking lot is already created",
			parkinglot: &ParkingLot{address: address, emptySlot: qheap.PriorityQueue{}, slots: slots, highestSlot: 0, capacity: 10},
			args: args{
				address:  address,
				capacity: 10,
			},
			want:    &ParkingLot{address: address, emptySlot: qheap.PriorityQueue{}, slots: slots, highestSlot: 0, capacity: 10},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.parkinglot.createParkingLot(tt.args.address, tt.args.capacity)
			if (err != nil) != tt.wantErr {
				t.Errorf("createParkingLot() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			compareParkingLot(t, tt.parkinglot, tt.want)
		})
	}
}

func TestGetNearestParkingSlot(t *testing.T) {
	address := "Marina Bay Sands"
	state := &ParkingLot{address: address, emptySlot: qheap.PriorityQueue{}, slots: generateParkingSlot(2), highestSlot: 0, capacity: 2}

	tests := []struct {
		name       string
		parkinglot *ParkingLot
		// fields  fields
		want    int
		wantErr bool
	}{
		{
			name:       "Empty parking lot with 2 available slots",
			parkinglot: state,
			want:       1,
			wantErr:    false,
		},
		{
			name:       "2 slots parking lot with 1 available slots",
			parkinglot: state,
			want:       2,
			wantErr:    false,
		},
		{
			name:       "Parking lot with unavailable slots",
			parkinglot: state,
			want:       0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parkinglot.getNearestParkingSlot()
			// fmt.Printf("test: %v, got: %v, err: %v\n", tt.name, got, err)

			if (err != nil) != tt.wantErr {
				t.Errorf("getNearestParkingSlot() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getNearestParkingSlot() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestPark(t *testing.T) {
	// Test data
	vehicle0 := &Vehicle{registrationNumber: "park KA-01-HH-2701", color: "Blue"}
	vehicle1 := &Vehicle{registrationNumber: "KA-01-HH-1234", color: "White"}
	vehicle2 := &Vehicle{registrationNumber: "KA-01-BB-0001", color: "Black"}

	emptySlot0 := qheap.PriorityQueue{}
	item1 := &qheap.Item{Value: 1}
	emptySlot1 := qheap.PriorityQueue{item1}

	slots := generateParkingSlot(2)
	address := "Marina Bay Sands"
	slots[0].vehicle = vehicle2
	slotAfterParkedByVehicle2 := slots[0] // slots[0] is slot marked with number 1
	// vehicle2 left slot 1, and then vehicle1 park at slot 1
	slots[0].vehicle = vehicle1
	slotAfterParkedByVehicle1 := slots[0]

	type args struct {
		registrationNumber string
		color              string
	}

	tests := []struct {
		name string
		// fields  fields
		parkinglot     *ParkingLot
		args           args
		wantSlot       *Slot
		wantErr        bool
		wantParkingLot *ParkingLot
	}{
		{
			name:       "ParkingLot is not created",
			parkinglot: &ParkingLot{},
			args: args{
				registrationNumber: vehicle1.registrationNumber,
				color:              vehicle1.color,
			},
			wantSlot:       nil,
			wantErr:        true,
			wantParkingLot: &ParkingLot{},
		},
		{
			name:       "Park vehicle into new slot",
			parkinglot: &ParkingLot{address: address, emptySlot: emptySlot0, slots: slots, highestSlot: 0, capacity: 2},
			args: args{
				registrationNumber: vehicle2.registrationNumber,
				color:              vehicle2.color,
			},
			wantSlot:       slotAfterParkedByVehicle2, // expected slotNumber = 1, vehicle2 with registrationNumber = KA-01-BB-0001
			wantErr:        false,
			wantParkingLot: &ParkingLot{address: address, emptySlot: emptySlot0, slots: slots, highestSlot: 1, capacity: 2},
		},
		{
			name:       "Park vehicle into a previously occupied but now free slot",
			parkinglot: &ParkingLot{address: address, emptySlot: emptySlot1, slots: slots, highestSlot: 1, capacity: 2},
			args: args{
				registrationNumber: vehicle1.registrationNumber,
				color:              vehicle1.color,
			},
			wantSlot:       slotAfterParkedByVehicle1, // expected slotNumber = 1, vehicle1 with registrationNumber = KA-01-HH-1234
			wantErr:        false,
			wantParkingLot: &ParkingLot{address: address, emptySlot: emptySlot0, slots: slots, highestSlot: 1, capacity: 2},
		},
		{
			name:       "Park car when parking lot is full",
			parkinglot: &ParkingLot{address: address, emptySlot: emptySlot0, slots: slots, highestSlot: 2, capacity: 2},
			args: args{
				registrationNumber: vehicle0.registrationNumber,
				color:              vehicle0.color,
			},
			wantSlot:       nil,
			wantErr:        true,
			wantParkingLot: &ParkingLot{address: address, emptySlot: emptySlot0, slots: slots, highestSlot: 2, capacity: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parkinglot.park(tt.args.registrationNumber, tt.args.color)
			// fmt.Printf("test: %v, got: %v, err: %v\n", tt.name, got, err)
			// if got != nil {
			// 	fmt.Printf("test: %v, gotSlot: %v, vehicle num: %v\n", tt.name, got.getParkingSlotNumber(), got.vehicle.registrationNumber)
			// }

			if (err != nil) != tt.wantErr {
				t.Errorf("park() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.wantSlot {
				t.Errorf("park() got = %v, wantSlot %v", got, tt.wantSlot)
			}

			compareParkingLot(t, tt.parkinglot, tt.wantParkingLot)
		})
	}
}
