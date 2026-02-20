package dtos

// CreateDriverInput contains the data for driver registration
type CreateDriverInput struct {
	DriverLicense string `json:"driverLicense" validate:"required,min=1"`
}
