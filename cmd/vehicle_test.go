package cmd

import (
	"reflect"
	"testing"
)

func TestCreateVehicle(t *testing.T) {
	type args struct {
		registrationNumber string
		color              string
	}
	tests := []struct {
		name string
		args args
		want *Vehicle
	}{
		{
			name: "Create vehicle",
			args: args{
				registrationNumber: "KA-01-HH-1234",
				color:              "White",
			},
			want: &Vehicle{
				registrationNumber: "KA-01-HH-1234",
				color:              "White",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createVehicle(tt.args.registrationNumber, tt.args.color)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createVehicle() got = %v, want %v", got, tt.want)
			}
		})
	}
}
