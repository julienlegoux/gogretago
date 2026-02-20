package entities

import "time"

// User represents the user profile domain entity
type User struct {
	ID           string
	RefID        int64
	FirstName    *string
	LastName     *string
	Phone        *string
	AuthRefID    int64
	AnonymizedAt *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// PublicUser extends User with email from the joined Auth record
type PublicUser struct {
	User
	Email string
}

// CreateUserData contains the data needed to create a new user profile
type CreateUserData struct {
	FirstName *string
	LastName  *string
	Phone     *string
	AuthRefID int64
}

// UpdateUserData contains partial update fields for a user profile
type UpdateUserData struct {
	FirstName *string
	LastName  *string
	Phone     *string
}
