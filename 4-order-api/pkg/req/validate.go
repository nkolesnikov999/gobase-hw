package req

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func IsValid[T any](payload T) error {
	validate := validator.New()
	// Accept 11 digits starting with 7 or 8. Adjust if needed.
	_ = validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		re := regexp.MustCompile(`^(7|8)\d{10}$`)
		return re.MatchString(phone)
	})
	err := validate.Struct(payload)
	return err
}
