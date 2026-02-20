package dtos

// CreateTripInput contains the data for creating a trip
type CreateTripInput struct {
	Kms           int    `json:"kms" validate:"required,gt=0"`
	Date          string `json:"date" validate:"required,min=1"`
	DepartureCity string `json:"departureCity" validate:"required,min=1"`
	ArrivalCity   string `json:"arrivalCity" validate:"required,min=1"`
	Seats         int    `json:"seats" validate:"required,gt=0"`
	CarID         string `json:"carId" validate:"required,min=1"`
}

// FindTripQuery contains the search query parameters
type FindTripQuery struct {
	DepartureCity *string `form:"departureCity"`
	ArrivalCity   *string `form:"arrivalCity"`
	Date          *string `form:"date"`
}
