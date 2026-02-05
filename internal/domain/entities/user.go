package entities

import "time"

// User represents the domain user entity
type User struct {
	ID        string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUserData contains the data needed to create a new user
type CreateUserData struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
}
