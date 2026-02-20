package dtos

// UpdateProfileInput contains the data for updating a user profile
type UpdateProfileInput struct {
	FirstName string `json:"firstName" validate:"required,min=1"`
	LastName  string `json:"lastName" validate:"required,min=1"`
	Phone     string `json:"phone" validate:"required,min=10"`
}
