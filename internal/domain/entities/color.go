package entities

// Color represents a car color domain entity
type Color struct {
	ID    string
	RefID int64
	Name  string
	Hex   string
}

// CreateColorData contains the data needed to create a new color
type CreateColorData struct {
	Name string
	Hex  string
}

// UpdateColorData contains partial update fields for a color
type UpdateColorData struct {
	Name *string
	Hex  *string
}
