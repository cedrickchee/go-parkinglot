package cmd

import (
	"reflect"
	"testing"
)

func compareSlot(t *testing.T, got *Slot, want *Slot) {
	if !reflect.DeepEqual(got.vehicle, want.vehicle) ||
		got.slotNumber != want.slotNumber {
		t.Errorf("parkVehicle() got = %v, want = %v", got, want)
	}
}

func TestSlotParkVehicle(t *testing.T) {
	type args struct {
		v *Vehicle
	}
	tests := []struct {
		name string
		slot *Slot
		args args
		want *Slot
	}{
		{
			name: "Park vehicle",
			slot: &Slot{
				slotNumber: 1,
			},
			args: args{
				v: &Vehicle{
					registrationNumber: "park KA-01-HH-1234",
					color:              "White",
				},
			},
			want: &Slot{
				vehicle: &Vehicle{
					registrationNumber: "park KA-01-HH-1234",
					color:              "White",
				},
				slotNumber: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.slot.parkVehicle(tt.args.v)
			compareSlot(t, tt.slot, tt.want)
		})
	}
}
