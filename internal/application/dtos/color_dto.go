package dtos

// CreateColorInput contains the data for creating a color
type CreateColorInput struct {
	Name string `json:"name" validate:"required,min=1"`
	Hex  string `json:"hex" validate:"required,hexcolor"`
}

// UpdateColorInput contains the data for updating a color
type UpdateColorInput struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Hex  *string `json:"hex,omitempty" validate:"omitempty,hexcolor"`
}
