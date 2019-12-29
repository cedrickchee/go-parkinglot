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

type fields struct {
	address    string
	vehicle0   *Vehicle
	vehicle1   *Vehicle
	vehicle2   *Vehicle
	slots      []*Slot
	item1      *qheap.Item
	emptySlot0 qheap.PriorityQueue
	emptySlot1 qheap.PriorityQueue
}

func genData() fields {
	data := fields{
		address:    "Marina Bay Sands",
		vehicle0:   &Vehicle{registrationNumber: "KA-01-HH-2701", color: "Blue"},
		vehicle1:   &Vehicle{registrationNumber: "KA-01-HH-1234", color: "White"},
		vehicle2:   &Vehicle{registrationNumber: "KA-01-BB-0001", color: "Black"},
		slots:      generateParkingSlot(10),
		item1:      &qheap.Item{Value: 1},
		emptySlot0: qheap.PriorityQueue{},
	}
	data.emptySlot1 = qheap.PriorityQueue{data.item1}

	return data
}

func compareParkingLot(t *testing.T, got *ParkingLot, want *ParkingLot) {
	if !reflect.DeepEqual(got.emptySlot, want.emptySlot) ||
		!reflect.DeepEqual(got.slots, want.slots) ||
		got.address != want.address ||
		got.highestSlot != want.highestSlot ||
		got.capacity != want.capacity {
		t.Errorf("ParkingLot got = %v, want = %v", got, want)
	}
}

func TestCreateParkingLot(t *testing.T) {
	data := genData()

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
				address:  data.address,
				capacity: 10,
			},
			want:    &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: data.slots, highestSlot: 0, capacity: 10},
			wantErr: false,
		},
		{
			name:       "Parking lot is already created",
			parkinglot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: data.slots, highestSlot: 0, capacity: 10},
			args: args{
				address:  data.address,
				capacity: 10,
			},
			want:    &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: data.slots, highestSlot: 0, capacity: 10},
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
	data := genData()
	state := &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: generateParkingSlot(2), highestSlot: 0, capacity: 2}

	tests := []struct {
		name       string
		parkinglot *ParkingLot
		want       int
		wantErr    bool
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
	data := genData()

	slots := generateParkingSlot(2)
	slots[0].vehicle = data.vehicle2
	slotAfterParkedByVehicle2 := slots[0] // slots[0] is slot marked with number 1
	// vehicle2 left slot 1, and then vehicle1 park at slot 1
	slots[0].vehicle = data.vehicle1
	slotAfterParkedByVehicle1 := slots[0]

	type args struct {
		registrationNumber string
		color              string
	}

	tests := []struct {
		name           string
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
				registrationNumber: data.vehicle1.registrationNumber,
				color:              data.vehicle1.color,
			},
			wantSlot:       nil,
			wantErr:        true,
			wantParkingLot: &ParkingLot{},
		},
		{
			name:       "Park vehicle into new slot",
			parkinglot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 0, capacity: 2},
			args: args{
				registrationNumber: data.vehicle2.registrationNumber,
				color:              data.vehicle2.color,
			},
			wantSlot:       slotAfterParkedByVehicle2, // expected slotNumber = 1, vehicle2 with registrationNumber = KA-01-BB-0001
			wantErr:        false,
			wantParkingLot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 1, capacity: 2},
		},
		{
			name:       "Park vehicle into a previously occupied but now free slot",
			parkinglot: &ParkingLot{address: data.address, emptySlot: data.emptySlot1, slots: slots, highestSlot: 1, capacity: 2},
			args: args{
				registrationNumber: data.vehicle1.registrationNumber,
				color:              data.vehicle1.color,
			},
			wantSlot:       slotAfterParkedByVehicle1, // expected slotNumber = 1, vehicle1 with registrationNumber = KA-01-HH-1234
			wantErr:        false,
			wantParkingLot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 1, capacity: 2},
		},
		{
			name:       "Park car when parking lot is full",
			parkinglot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 2, capacity: 2},
			args: args{
				registrationNumber: data.vehicle0.registrationNumber,
				color:              data.vehicle0.color,
			},
			wantSlot:       nil,
			wantErr:        true,
			wantParkingLot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 2, capacity: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parkinglot.park(tt.args.registrationNumber, tt.args.color)

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

func TestLeave(t *testing.T) {
	data := genData()

	slots := data.slots
	slots[0].vehicle = data.vehicle0

	type args struct {
		slotNumber int
	}

	tests := []struct {
		name           string
		parkinglot     *ParkingLot
		args           args
		wantErr        bool
		wantParkingLot *ParkingLot
	}{
		{
			name:           "ParkingLot is not created",
			parkinglot:     &ParkingLot{},
			args:           args{slotNumber: 1},
			wantErr:        true,
			wantParkingLot: &ParkingLot{},
		},
		{
			name:           "Leave existing vehicle",
			parkinglot:     &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 2, capacity: 10},
			args:           args{slotNumber: 1},
			wantErr:        false,
			wantParkingLot: &ParkingLot{address: data.address, emptySlot: data.emptySlot1, slots: slots, highestSlot: 2, capacity: 10},
		},
		{
			name:           "Leave non-existent vehicle",
			parkinglot:     &ParkingLot{address: data.address, emptySlot: data.emptySlot1, slots: slots, highestSlot: 2, capacity: 10},
			args:           args{slotNumber: 2},
			wantErr:        true,
			wantParkingLot: &ParkingLot{address: data.address, emptySlot: data.emptySlot1, slots: slots, highestSlot: 2, capacity: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.parkinglot.leave(tt.args.slotNumber)

			if (err != nil) != tt.wantErr {
				t.Errorf("leave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			compareParkingLot(t, tt.parkinglot, tt.wantParkingLot)
		})
	}
}

func TestGetStatus(t *testing.T) {
	data := genData()

	slots := data.slots
	slots[0].vehicle = data.vehicle1
	slots[1].vehicle = data.vehicle2

	tests := []struct {
		name       string
		parkinglot *ParkingLot
		want       []*Slot
	}{
		{
			name:       "Parking lot is not created",
			parkinglot: &ParkingLot{},
			want:       nil,
		},
		{
			name:       "Parking lot is empty",
			parkinglot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 0, capacity: 10},
			want:       nil,
		},
		{
			name:       "Parking lot with vehicles",
			parkinglot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 2, capacity: 10},
			want: []*Slot{
				{slotNumber: 1, vehicle: data.vehicle1},
				{slotNumber: 2, vehicle: data.vehicle2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.parkinglot.getStatus(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVehiclesByColor(t *testing.T) {
	data := genData()

	slots := data.slots
	slots[0].vehicle = data.vehicle1

	type args struct {
		color string
	}

	tests := []struct {
		name         string
		parkinglot   *ParkingLot
		args         args
		wantSlot     []int
		wantRegisNum []string
		wantErr      bool
	}{
		{
			name:         "Parking lot is not created",
			parkinglot:   &ParkingLot{},
			args:         args{color: "White"},
			wantSlot:     nil,
			wantRegisNum: nil,
			wantErr:      true,
		},
		{
			name:         "A vehicle is parked and the color is White",
			parkinglot:   &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 1, capacity: 10},
			args:         args{color: "White"},
			wantSlot:     []int{1},
			wantRegisNum: []string{"KA-01-HH-1234"},
			wantErr:      false,
		},
		{
			name:         "A vehicle is not parked with the requested color",
			parkinglot:   &ParkingLot{address: data.address, emptySlot: data.emptySlot1, slots: slots, highestSlot: 2, capacity: 10},
			args:         args{color: "Black"},
			wantSlot:     nil,
			wantRegisNum: nil,
			wantErr:      true,
		},
		{
			name:         "Parking lot is empty",
			parkinglot:   &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 0, capacity: 10},
			args:         args{color: "White"},
			wantSlot:     nil,
			wantRegisNum: nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSlots, gotRegisNumbers, err := tt.parkinglot.getVehiclesByColor(tt.args.color)

			if (err != nil) != tt.wantErr {
				t.Errorf("getVehiclesByColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSlots, tt.wantSlot) {
				t.Errorf("getVehiclesByColor() gotSlots = %v, want %v", gotSlots, tt.wantSlot)
			}
			if !reflect.DeepEqual(gotRegisNumbers, tt.wantRegisNum) {
				t.Errorf("getVehiclesByColor() gotRegisNumbers = %v, want %v", gotRegisNumbers, tt.wantRegisNum)
			}
		})
	}
}

func TestGetVehicleByRegistrationNumber(t *testing.T) {
	data := genData()

	slots := data.slots
	slots[0].vehicle = data.vehicle1

	type args struct {
		registrationNumber string
	}

	tests := []struct {
		name       string
		parkinglot *ParkingLot
		args       args
		want       int
		wantErr    bool
	}{
		{
			name:       "Parking lot is not created",
			parkinglot: &ParkingLot{},
			args:       args{registrationNumber: "KA-01-HH-1234"},
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Parking lot has vehicle of given registration number",
			parkinglot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 1, capacity: 10},
			args:       args{registrationNumber: "KA-01-HH-1234"},
			want:       1,
			wantErr:    false,
		},
		{
			name:       "Parking lot don't have vehicle of given registration number",
			parkinglot: &ParkingLot{address: data.address, emptySlot: data.emptySlot1, slots: slots, highestSlot: 2, capacity: 10},
			args:       args{registrationNumber: "KA-01-BB-0001"},
			want:       0,
			wantErr:    true,
		},
		{
			name:       "Parking lot is empty",
			parkinglot: &ParkingLot{address: data.address, emptySlot: data.emptySlot0, slots: slots, highestSlot: 0, capacity: 10},
			args:       args{registrationNumber: "KA-01-HH-1234"},
			want:       0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.parkinglot.getVehicleByRegistrationNumber(tt.args.registrationNumber)

			if (err != nil) != tt.wantErr {
				t.Errorf("getVehicleByRegistrationNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getVehicleByRegistrationNumber() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
