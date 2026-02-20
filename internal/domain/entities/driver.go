package entities

import "time"

// Driver represents a driver domain entity
type Driver struct {
	ID            string
	RefID         int64
	DriverLicense string
	UserRefID     int64
	AnonymizedAt  *time.Time
}

// CreateDriverData contains the data needed to create a new driver
type CreateDriverData struct {
	DriverLicense string
	UserRefID     int64
}
