package dtos

// CreateBrandInput contains the data for creating a brand
type CreateBrandInput struct {
	Name string `json:"name" validate:"required,min=1"`
}
