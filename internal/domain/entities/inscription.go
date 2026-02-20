package entities

import "time"

// Inscription represents a passenger booking domain entity
type Inscription struct {
	ID        string
	RefID     int64
	CreatedAt time.Time
	UserRefID int64
	TripRefID int64
	Status    string
}

// CreateInscriptionData contains the data needed to create a new inscription
type CreateInscriptionData struct {
	UserRefID int64
	TripRefID int64
}
