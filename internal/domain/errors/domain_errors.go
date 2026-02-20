package errors

import "fmt"

// DomainError is the base error type for domain-level errors
type DomainError struct {
	Message string
	Code    string
}

func (e *DomainError) Error() string {
	return e.Message
}

// --- Error registry: maps error codes to HTTP status codes ---

var errorHTTPStatus = map[string]int{
	"USER_ALREADY_EXISTS":   409,
	"INVALID_CREDENTIALS":   401,
	"USER_NOT_FOUND":        404,
	"BRAND_NOT_FOUND":       404,
	"CITY_NOT_FOUND":        404,
	"CAR_NOT_FOUND":         404,
	"CAR_ALREADY_EXISTS":    409,
	"DRIVER_NOT_FOUND":      404,
	"DRIVER_ALREADY_EXISTS":  409,
	"TRIP_NOT_FOUND":        404,
	"INSCRIPTION_NOT_FOUND": 404,
	"ALREADY_INSCRIBED":     409,
	"NO_SEATS_AVAILABLE":    400,
	"COLOR_NOT_FOUND":       404,
	"COLOR_ALREADY_EXISTS":  409,
	"FORBIDDEN":             403,
	"UNAUTHORIZED":          401,
	"TOKEN_EXPIRED":         401,
	"TOKEN_INVALID":         401,
	"TOKEN_MALFORMED":       400,
	"VALIDATION_ERROR":      400,
	"RELATION_CONSTRAINT":   409,
	"INTERNAL_ERROR":        500,
}

// GetHTTPStatus returns the HTTP status code for a given error code
func GetHTTPStatus(code string) int {
	if status, ok := errorHTTPStatus[code]; ok {
		return status
	}
	return 500
}

// --- Concrete domain error types ---

type UserAlreadyExistsError struct{ DomainError }

func NewUserAlreadyExistsError(email string) *UserAlreadyExistsError {
	return &UserAlreadyExistsError{DomainError{
		Message: fmt.Sprintf("A user with email %q already exists", email),
		Code:    "USER_ALREADY_EXISTS",
	}}
}

type InvalidCredentialsError struct{ DomainError }

func NewInvalidCredentialsError() *InvalidCredentialsError {
	return &InvalidCredentialsError{DomainError{
		Message: "Invalid email or password",
		Code:    "INVALID_CREDENTIALS",
	}}
}

type UserNotFoundError struct{ DomainError }

func NewUserNotFoundError(identifier string) *UserNotFoundError {
	return &UserNotFoundError{DomainError{
		Message: fmt.Sprintf("User not found: %s", identifier),
		Code:    "USER_NOT_FOUND",
	}}
}

type BrandNotFoundError struct{ DomainError }

func NewBrandNotFoundError(identifier string) *BrandNotFoundError {
	return &BrandNotFoundError{DomainError{
		Message: fmt.Sprintf("Brand not found: %s", identifier),
		Code:    "BRAND_NOT_FOUND",
	}}
}

type CityNotFoundError struct{ DomainError }

func NewCityNotFoundError(identifier string) *CityNotFoundError {
	return &CityNotFoundError{DomainError{
		Message: fmt.Sprintf("City not found: %s", identifier),
		Code:    "CITY_NOT_FOUND",
	}}
}

type CarNotFoundError struct{ DomainError }

func NewCarNotFoundError(identifier string) *CarNotFoundError {
	return &CarNotFoundError{DomainError{
		Message: fmt.Sprintf("Car not found: %s", identifier),
		Code:    "CAR_NOT_FOUND",
	}}
}

type CarAlreadyExistsError struct{ DomainError }

func NewCarAlreadyExistsError(licensePlate string) *CarAlreadyExistsError {
	return &CarAlreadyExistsError{DomainError{
		Message: fmt.Sprintf("A car with license plate %q already exists", licensePlate),
		Code:    "CAR_ALREADY_EXISTS",
	}}
}

type DriverNotFoundError struct{ DomainError }

func NewDriverNotFoundError(identifier string) *DriverNotFoundError {
	return &DriverNotFoundError{DomainError{
		Message: fmt.Sprintf("Driver not found: %s", identifier),
		Code:    "DRIVER_NOT_FOUND",
	}}
}

type DriverAlreadyExistsError struct{ DomainError }

func NewDriverAlreadyExistsError(userId string) *DriverAlreadyExistsError {
	return &DriverAlreadyExistsError{DomainError{
		Message: fmt.Sprintf("A driver already exists for user %q", userId),
		Code:    "DRIVER_ALREADY_EXISTS",
	}}
}

type TripNotFoundError struct{ DomainError }

func NewTripNotFoundError(identifier string) *TripNotFoundError {
	return &TripNotFoundError{DomainError{
		Message: fmt.Sprintf("Trip not found: %s", identifier),
		Code:    "TRIP_NOT_FOUND",
	}}
}

type InscriptionNotFoundError struct{ DomainError }

func NewInscriptionNotFoundError(identifier string) *InscriptionNotFoundError {
	return &InscriptionNotFoundError{DomainError{
		Message: fmt.Sprintf("Inscription not found: %s", identifier),
		Code:    "INSCRIPTION_NOT_FOUND",
	}}
}

type AlreadyInscribedError struct{ DomainError }

func NewAlreadyInscribedError(userId, tripId string) *AlreadyInscribedError {
	return &AlreadyInscribedError{DomainError{
		Message: fmt.Sprintf("User %s is already inscribed to trip %s", userId, tripId),
		Code:    "ALREADY_INSCRIBED",
	}}
}

type NoSeatsAvailableError struct{ DomainError }

func NewNoSeatsAvailableError(tripId string) *NoSeatsAvailableError {
	return &NoSeatsAvailableError{DomainError{
		Message: fmt.Sprintf("No seats available on trip %s", tripId),
		Code:    "NO_SEATS_AVAILABLE",
	}}
}

type ColorNotFoundError struct{ DomainError }

func NewColorNotFoundError(id string) *ColorNotFoundError {
	return &ColorNotFoundError{DomainError{
		Message: fmt.Sprintf("Color not found: %s", id),
		Code:    "COLOR_NOT_FOUND",
	}}
}

type ColorAlreadyExistsError struct{ DomainError }

func NewColorAlreadyExistsError(name string) *ColorAlreadyExistsError {
	return &ColorAlreadyExistsError{DomainError{
		Message: fmt.Sprintf("Color already exists: %s", name),
		Code:    "COLOR_ALREADY_EXISTS",
	}}
}

type ForbiddenError struct{ DomainError }

func NewForbiddenError(resource, id string) *ForbiddenError {
	return &ForbiddenError{DomainError{
		Message: fmt.Sprintf("You do not have permission to access %s: %s", resource, id),
		Code:    "FORBIDDEN",
	}}
}
