package cmd

import (
	"fmt"
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
			fmt.Printf("test: %v, got: %v, err: %v\n", tt.name, got, err)

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
