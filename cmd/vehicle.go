package cmd

type Vehicle struct {
	registrationNumber string
	color              string
}

// Create a new vehicle
func createVehicle(registrationNumber, color string) *Vehicle {
	return &Vehicle{registrationNumber, color}
}

// Returns vehicle registration number
func (v *Vehicle) getNumber() string {
	return v.registrationNumber
}

// Returns vehicle color
func (v *Vehicle) getColor() string {
	return v.color
}
