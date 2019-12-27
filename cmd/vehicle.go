package cmd

type Vehicle struct {
	registrationNumber string
	color              string
}

// Create a new vehicle
func createVehicle(registrationNumber, color string) *Vehicle {
	return &Vehicle{registrationNumber, color}
}
