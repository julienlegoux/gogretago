package entities

// CityTrip represents the association between a city and a trip
type CityTrip struct {
	TripRefID int64
	CityRefID int64
	Type      string // "DEPARTURE" or "ARRIVAL"
}
