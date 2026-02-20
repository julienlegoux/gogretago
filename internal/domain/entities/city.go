package entities

// City represents a city domain entity
type City struct {
	ID       string
	RefID    int64
	CityName string
	Zipcode  string
}

// CreateCityData contains the data needed to create a new city
type CreateCityData struct {
	CityName string
	Zipcode  string
}
