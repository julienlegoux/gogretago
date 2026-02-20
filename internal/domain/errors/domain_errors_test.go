package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHTTPStatus_AllKnownCodes(t *testing.T) {
	expected := map[string]int{
		"USER_ALREADY_EXISTS":   409,
		"INVALID_CREDENTIALS":   401,
		"USER_NOT_FOUND":        404,
		"BRAND_NOT_FOUND":       404,
		"CITY_NOT_FOUND":        404,
		"CAR_NOT_FOUND":         404,
		"CAR_ALREADY_EXISTS":    409,
		"DRIVER_NOT_FOUND":      404,
		"DRIVER_ALREADY_EXISTS": 409,
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

	for code, status := range expected {
		t.Run(code, func(t *testing.T) {
			assert.Equal(t, status, GetHTTPStatus(code))
		})
	}
}

func TestGetHTTPStatus_UnknownCode(t *testing.T) {
	assert.Equal(t, 500, GetHTTPStatus("UNKNOWN"))
}

func TestDomainError_ErrorMethod(t *testing.T) {
	err := &DomainError{Message: "something went wrong", Code: "TEST"}
	assert.Equal(t, "something went wrong", err.Error())
}

func TestNewUserAlreadyExistsError(t *testing.T) {
	err := NewUserAlreadyExistsError("test@example.com")
	assert.Equal(t, "USER_ALREADY_EXISTS", err.Code)
	assert.Contains(t, err.Message, "test@example.com")
}

func TestNewInvalidCredentialsError(t *testing.T) {
	err := NewInvalidCredentialsError()
	assert.Equal(t, "INVALID_CREDENTIALS", err.Code)
	assert.Equal(t, "Invalid email or password", err.Message)
}

func TestNewUserNotFoundError(t *testing.T) {
	err := NewUserNotFoundError("abc-123")
	assert.Equal(t, "USER_NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "abc-123")
}

func TestNewBrandNotFoundError(t *testing.T) {
	err := NewBrandNotFoundError("brand-1")
	assert.Equal(t, "BRAND_NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "brand-1")
}

func TestNewCityNotFoundError(t *testing.T) {
	err := NewCityNotFoundError("city-1")
	assert.Equal(t, "CITY_NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "city-1")
}

func TestNewCarNotFoundError(t *testing.T) {
	err := NewCarNotFoundError("car-1")
	assert.Equal(t, "CAR_NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "car-1")
}

func TestNewCarAlreadyExistsError(t *testing.T) {
	err := NewCarAlreadyExistsError("ABC-123")
	assert.Equal(t, "CAR_ALREADY_EXISTS", err.Code)
	assert.Contains(t, err.Message, "ABC-123")
}

func TestNewDriverNotFoundError(t *testing.T) {
	err := NewDriverNotFoundError("drv-1")
	assert.Equal(t, "DRIVER_NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "drv-1")
}

func TestNewDriverAlreadyExistsError(t *testing.T) {
	err := NewDriverAlreadyExistsError("user-1")
	assert.Equal(t, "DRIVER_ALREADY_EXISTS", err.Code)
	assert.Contains(t, err.Message, "user-1")
}

func TestNewTripNotFoundError(t *testing.T) {
	err := NewTripNotFoundError("trip-1")
	assert.Equal(t, "TRIP_NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "trip-1")
}

func TestNewInscriptionNotFoundError(t *testing.T) {
	err := NewInscriptionNotFoundError("insc-1")
	assert.Equal(t, "INSCRIPTION_NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "insc-1")
}

func TestNewAlreadyInscribedError(t *testing.T) {
	err := NewAlreadyInscribedError("user-1", "trip-1")
	assert.Equal(t, "ALREADY_INSCRIBED", err.Code)
	assert.Contains(t, err.Message, "user-1")
	assert.Contains(t, err.Message, "trip-1")
}

func TestNewNoSeatsAvailableError(t *testing.T) {
	err := NewNoSeatsAvailableError("trip-1")
	assert.Equal(t, "NO_SEATS_AVAILABLE", err.Code)
	assert.Contains(t, err.Message, "trip-1")
}

func TestNewColorNotFoundError(t *testing.T) {
	err := NewColorNotFoundError("color-1")
	assert.Equal(t, "COLOR_NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "color-1")
}

func TestNewColorAlreadyExistsError(t *testing.T) {
	err := NewColorAlreadyExistsError("Red")
	assert.Equal(t, "COLOR_ALREADY_EXISTS", err.Code)
	assert.Contains(t, err.Message, "Red")
}

func TestNewForbiddenError(t *testing.T) {
	err := NewForbiddenError("trip", "trip-1")
	assert.Equal(t, "FORBIDDEN", err.Code)
	assert.Contains(t, err.Message, "trip")
	assert.Contains(t, err.Message, "trip-1")
}

func TestDomainErrors_ImplementErrorInterface(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"UserAlreadyExistsError", NewUserAlreadyExistsError("a@b.com")},
		{"InvalidCredentialsError", NewInvalidCredentialsError()},
		{"UserNotFoundError", NewUserNotFoundError("1")},
		{"BrandNotFoundError", NewBrandNotFoundError("1")},
		{"CityNotFoundError", NewCityNotFoundError("1")},
		{"CarNotFoundError", NewCarNotFoundError("1")},
		{"CarAlreadyExistsError", NewCarAlreadyExistsError("ABC")},
		{"DriverNotFoundError", NewDriverNotFoundError("1")},
		{"DriverAlreadyExistsError", NewDriverAlreadyExistsError("1")},
		{"TripNotFoundError", NewTripNotFoundError("1")},
		{"InscriptionNotFoundError", NewInscriptionNotFoundError("1")},
		{"AlreadyInscribedError", NewAlreadyInscribedError("1", "2")},
		{"NoSeatsAvailableError", NewNoSeatsAvailableError("1")},
		{"ColorNotFoundError", NewColorNotFoundError("1")},
		{"ColorAlreadyExistsError", NewColorAlreadyExistsError("Red")},
		{"ForbiddenError", NewForbiddenError("res", "1")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Implements(t, (*error)(nil), tt.err)
			assert.NotEmpty(t, tt.err.Error())
		})
	}
}

func TestDomainErrors_EmbedDomainError(t *testing.T) {
	// Typed domain errors embed DomainError by value.
	// In Go, errors.As does NOT match embedded value types (only pointers or Unwrap chains).
	// Instead, verify the embedded DomainError fields are accessible directly
	// and that a plain *DomainError works with errors.As.
	tests := []struct {
		name    string
		err     error
		code    string
	}{
		{"UserAlreadyExistsError", NewUserAlreadyExistsError("a@b.com"), "USER_ALREADY_EXISTS"},
		{"InvalidCredentialsError", NewInvalidCredentialsError(), "INVALID_CREDENTIALS"},
		{"UserNotFoundError", NewUserNotFoundError("1"), "USER_NOT_FOUND"},
		{"BrandNotFoundError", NewBrandNotFoundError("1"), "BRAND_NOT_FOUND"},
		{"CityNotFoundError", NewCityNotFoundError("1"), "CITY_NOT_FOUND"},
		{"CarNotFoundError", NewCarNotFoundError("1"), "CAR_NOT_FOUND"},
		{"CarAlreadyExistsError", NewCarAlreadyExistsError("ABC"), "CAR_ALREADY_EXISTS"},
		{"DriverNotFoundError", NewDriverNotFoundError("1"), "DRIVER_NOT_FOUND"},
		{"DriverAlreadyExistsError", NewDriverAlreadyExistsError("1"), "DRIVER_ALREADY_EXISTS"},
		{"TripNotFoundError", NewTripNotFoundError("1"), "TRIP_NOT_FOUND"},
		{"InscriptionNotFoundError", NewInscriptionNotFoundError("1"), "INSCRIPTION_NOT_FOUND"},
		{"AlreadyInscribedError", NewAlreadyInscribedError("1", "2"), "ALREADY_INSCRIBED"},
		{"NoSeatsAvailableError", NewNoSeatsAvailableError("1"), "NO_SEATS_AVAILABLE"},
		{"ColorNotFoundError", NewColorNotFoundError("1"), "COLOR_NOT_FOUND"},
		{"ColorAlreadyExistsError", NewColorAlreadyExistsError("Red"), "COLOR_ALREADY_EXISTS"},
		{"ForbiddenError", NewForbiddenError("res", "1"), "FORBIDDEN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.err.Error())
			// Extract the embedded DomainError.Code field via reflection
			val := reflect.ValueOf(tt.err)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			codeField := val.FieldByName("Code")
			assert.True(t, codeField.IsValid(), "expected Code field on %T", tt.err)
			assert.Equal(t, tt.code, codeField.String())
		})
	}
}

func TestDomainError_WorksWithErrorsAs(t *testing.T) {
	// A plain *DomainError should be matchable via errors.As
	plainErr := &DomainError{Code: "USER_NOT_FOUND", Message: "User not found: 1"}
	var target *DomainError
	assert.True(t, errors.As(plainErr, &target))
	assert.Equal(t, "USER_NOT_FOUND", target.Code)
	assert.Equal(t, "User not found: 1", target.Message)
}

func TestConcreteErrors_ErrorsAs_DoesNotMatchDomainError(t *testing.T) {
	// Concrete typed errors embed DomainError by VALUE, so errors.As
	// with a *DomainError target does NOT match them. This is important
	// to document as expected behavior since the error handler relies on
	// receiving plain *DomainError values.
	concreteErrors := []error{
		NewUserAlreadyExistsError("a@b.com"),
		NewInvalidCredentialsError(),
		NewUserNotFoundError("1"),
		NewCarNotFoundError("1"),
		NewTripNotFoundError("1"),
	}

	for _, err := range concreteErrors {
		var target *DomainError
		assert.False(t, errors.As(err, &target),
			"errors.As should NOT match %T to *DomainError (embedded by value)", err)
	}
}

func TestConcreteErrors_ErrorsAs_MatchOwnType(t *testing.T) {
	// Each concrete error type should be matchable via errors.As to its own type
	t.Run("UserAlreadyExistsError", func(t *testing.T) {
		err := NewUserAlreadyExistsError("a@b.com")
		var target *UserAlreadyExistsError
		assert.True(t, errors.As(err, &target))
		assert.Equal(t, "USER_ALREADY_EXISTS", target.Code)
	})
	t.Run("UserNotFoundError", func(t *testing.T) {
		err := NewUserNotFoundError("1")
		var target *UserNotFoundError
		assert.True(t, errors.As(err, &target))
		assert.Equal(t, "USER_NOT_FOUND", target.Code)
	})
	t.Run("CarNotFoundError", func(t *testing.T) {
		err := NewCarNotFoundError("1")
		var target *CarNotFoundError
		assert.True(t, errors.As(err, &target))
		assert.Equal(t, "CAR_NOT_FOUND", target.Code)
	})
	t.Run("ForbiddenError", func(t *testing.T) {
		err := NewForbiddenError("res", "1")
		var target *ForbiddenError
		assert.True(t, errors.As(err, &target))
		assert.Equal(t, "FORBIDDEN", target.Code)
	})
	t.Run("ColorAlreadyExistsError", func(t *testing.T) {
		err := NewColorAlreadyExistsError("Red")
		var target *ColorAlreadyExistsError
		assert.True(t, errors.As(err, &target))
		assert.Equal(t, "COLOR_ALREADY_EXISTS", target.Code)
	})
}
