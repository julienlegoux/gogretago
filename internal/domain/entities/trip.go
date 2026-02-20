package entities

import "time"

// Trip represents a carpooling trip domain entity
type Trip struct {
	ID          string
	RefID       int64
	DateTrip    time.Time
	Kms         int
	Seats       int
	DriverRefID int64
	CarRefID    int64
}

// CreateTripData contains the data needed to create a new trip
type CreateTripData struct {
	DateTrip    time.Time
	Kms         int
	Seats       int
	DriverRefID int64
	CarRefID    int64
	CityRefIDs  []int64 // [departureRefID, arrivalRefID]
}

// TripFilters contains optional filters for searching trips
type TripFilters struct {
	DepartureCity *string
	ArrivalCity   *string
	Date          *time.Time
}
