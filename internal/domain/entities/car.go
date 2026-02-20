package entities

// Car represents a car domain entity
type Car struct {
	ID           string
	RefID        int64
	LicensePlate string
	ModelRefID   int64
	DriverRefID  int64
}

// CreateCarData contains the data needed to create a new car
type CreateCarData struct {
	LicensePlate string
	ModelRefID   int64
	DriverRefID  int64
}

// UpdateCarData contains partial update fields for a car
type UpdateCarData struct {
	LicensePlate *string
	ModelRefID   *int64
}
