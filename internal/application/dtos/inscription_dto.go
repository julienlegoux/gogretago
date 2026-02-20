package dtos

// CreateInscriptionInput contains the data for booking a trip
type CreateInscriptionInput struct {
	TripID string `json:"tripId" validate:"required,min=1"`
}
