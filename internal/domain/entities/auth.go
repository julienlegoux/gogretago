package entities

import "time"

// Auth represents the authentication domain entity
type Auth struct {
	ID           string
	RefID        int64
	Email        string
	Password     string
	Role         string
	AnonymizedAt *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CreateAuthData contains the data needed to create a new auth record
type CreateAuthData struct {
	Email    string
	Password string
}
