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

// UserAlreadyExistsError is returned when attempting to register with an existing email
type UserAlreadyExistsError struct {
	DomainError
}

func NewUserAlreadyExistsError(email string) *UserAlreadyExistsError {
	return &UserAlreadyExistsError{
		DomainError: DomainError{
			Message: fmt.Sprintf("A user with email \"%s\" already exists", email),
			Code:    "USER_ALREADY_EXISTS",
		},
	}
}

// InvalidCredentialsError is returned when login credentials are incorrect
type InvalidCredentialsError struct {
	DomainError
}

func NewInvalidCredentialsError() *InvalidCredentialsError {
	return &InvalidCredentialsError{
		DomainError: DomainError{
			Message: "Invalid email or password",
			Code:    "INVALID_CREDENTIALS",
		},
	}
}

// UserNotFoundError is returned when a user lookup fails
type UserNotFoundError struct {
	DomainError
}

func NewUserNotFoundError(identifier string) *UserNotFoundError {
	return &UserNotFoundError{
		DomainError: DomainError{
			Message: fmt.Sprintf("User not found: %s", identifier),
			Code:    "USER_NOT_FOUND",
		},
	}
}
