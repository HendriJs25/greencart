package error

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

func ErrValidationResponse(err error) (validationResponse []ValidationResponse) {
	var fieldErrors validator.ValidationErrors

	if errors.As(err, &fieldErrors) {
		for _, err := range fieldErrors {
			fieldJa := getFieldNameJa(err.StructField())
			paramJa := getFieldNameJa(err.Param())
			errValidator, ok := ErrValidator[err.Tag()]
			if ok {
				count := strings.Count(errValidator, "%s")
				if count == 1 {
					validationResponse = append(validationResponse, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf(errValidator, fieldJa),
					})
				} else {
					validationResponse = append(validationResponse, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf(errValidator, fieldJa, paramJa),
					})
				}
			} else {
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%sの入力内容が正しくありません", fieldJa),
				})
			}
		}
	}
	return validationResponse
}

func getFieldNameJa(field string) string {
	if name, ok := FieldNameJa[field]; ok {
		return name
	}
	return field
}
