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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.parkinglot.createParkingLot(tt.args.address, tt.args.capacity)
			if (err != nil) != tt.wantErr {
				t.Errorf("createParkingLot() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.parkinglot, tt.want) {
				t.Errorf("createParkingLot() got = %v, want = %v", tt.parkinglot, tt.want)
			}
		})
	}
}
