package dtos

// CreateCityInput contains the data for creating a city
type CreateCityInput struct {
	CityName string `json:"cityName" validate:"required,min=1"`
	Zipcode  string `json:"zipcode" validate:"required"`
}
