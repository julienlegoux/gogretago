package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom password validation
	validate.RegisterValidation("password", validatePassword)

	// Register custom French phone validation
	validate.RegisterValidation("frenchphone", validateFrenchPhone)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}

// validatePassword checks password requirements:
// - At least one lowercase letter
// - At least one uppercase letter
// - At least one digit
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	return hasLower && hasUpper && hasDigit
}

// validateFrenchPhone validates French phone number format
func validateFrenchPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	re := regexp.MustCompile(`^(\+33|0)[1-9](\d{2}){4}$`)
	return re.MatchString(phone)
}

// ValidationError represents a validation error with field details
type ValidationErrors struct {
	Errors map[string][]string
}

// FormatValidationErrors converts validator errors to a map
func FormatValidationErrors(err error) map[string][]string {
	errors := make(map[string][]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			var message string

			switch e.Tag() {
			case "required":
				message = field + " is required"
			case "email":
				message = "Invalid email format"
			case "min":
				message = field + " must be at least " + e.Param() + " characters"
			case "eqfield":
				message = "Passwords do not match"
			case "password":
				message = "Password must contain at least one lowercase, one uppercase, and one number"
			case "frenchphone":
				message = "Invalid phone number format"
			default:
				message = field + " is invalid"
			}

			errors[field] = append(errors[field], message)
		}
	}

	return errors
}
