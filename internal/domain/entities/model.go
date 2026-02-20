package entities

// VehicleModel represents a car model domain entity
type VehicleModel struct {
	ID         string
	RefID      int64
	Name       string
	BrandRefID int64
}

// CreateModelData contains the data needed to create a new car model
type CreateModelData struct {
	Name       string
	BrandRefID int64
}
