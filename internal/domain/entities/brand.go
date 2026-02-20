package entities

// Brand represents a car brand (manufacturer) domain entity
type Brand struct {
	ID    string
	RefID int64
	Name  string
}

// CreateBrandData contains the data needed to create a new brand
type CreateBrandData struct {
	Name string
}
