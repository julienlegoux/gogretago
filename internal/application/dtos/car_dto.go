package dtos

// CreateCarInput contains the data for creating a car
type CreateCarInput struct {
	Model        string `json:"model" validate:"required,min=1"`
	BrandID      string `json:"brandId" validate:"required,min=1"`
	LicensePlate string `json:"licensePlate" validate:"required,min=1"`
}

// UpdateCarInput contains the data for a full car update (PUT)
type UpdateCarInput struct {
	Model        string `json:"model" validate:"required,min=1"`
	BrandID      string `json:"brandId" validate:"required,min=1"`
	LicensePlate string `json:"licensePlate" validate:"required,min=1"`
}

// PatchCarInput contains the data for a partial car update (PATCH)
type PatchCarInput struct {
	Model        *string `json:"model,omitempty" validate:"omitempty,min=1"`
	BrandID      *string `json:"brandId,omitempty" validate:"omitempty,min=1"`
	LicensePlate *string `json:"licensePlate,omitempty" validate:"omitempty,min=1"`
}
