package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	lowerRegex  = regexp.MustCompile(`[a-z]`)
	upperRegex  = regexp.MustCompile(`[A-Z]`)
	numberRegex = regexp.MustCompile(`[0-9]`)
)

func New() (*validator.Validate, error) {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := RegisterCustomValidators(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func RegisterCustomValidators(v *validator.Validate) error {
	validations := map[string]validator.Func{
		"validpassword": validPassword,
		"notblank":      notBlank,
	}

	for tag, fn := range validations {
		if err := v.RegisterValidation(tag, fn); err != nil {
			return fmt.Errorf("failed to register custom validator %q :%w", tag, err)
		}
	}
	return nil
}

func validPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	return len(password) >= 8 &&
		lowerRegex.MatchString(password) &&
		upperRegex.MatchString(password) &&
		numberRegex.MatchString(password)
}

func notBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}
