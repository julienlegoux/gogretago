package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("hexcolor", validateHexColor)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	return hasLower && hasUpper && hasDigit
}

func validateHexColor(fl validator.FieldLevel) bool {
	hex := fl.Field().String()
	re := regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
	return re.MatchString(hex)
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
			case "gt":
				message = field + " must be greater than " + e.Param()
			case "eqfield":
				message = "Passwords do not match"
			case "password":
				message = "Password must contain at least one lowercase, one uppercase, and one number"
			case "hexcolor":
				message = "Must be a valid hex color (#RRGGBB)"
			default:
				message = field + " is invalid"
			}

			errors[field] = append(errors[field], message)
		}
	}

	return errors
}
