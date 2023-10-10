package validate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidateError struct {
	Field  string `json:"field"`
	Issue  string `json:"issue"`
	Option string `json:"option"`
	Input  string `json:"input"`
}

type BadRequestErrorResponse struct {
	Errors []ValidateError `json:"errors"`
	Url    string          `json:"url"`
	Method string          `json:"method"`
}

func MakeBadRequestErrorResponse(c echo.Context, err validator.ValidationErrors) BadRequestErrorResponse {

	var errors []ValidateError

	for _, err := range err {

		fieldName := err.Field()

		paramName := strings.Replace(fieldName, fieldName[0:1], strings.ToLower(fieldName[0:1]), 1)
		var input string
		switch v := err.Value().(type) {
		case int:
			input = strconv.Itoa(v)
		case string:
			input = v
		case []int:
			input = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v)), ","), "[]")

		}

		validateError := ValidateError{
			Field:  paramName,
			Issue:  err.Tag(),
			Option: err.Param(),
			Input:  input,
		}

		errors = append(errors, validateError)

	}
	return BadRequestErrorResponse{
		Errors: errors,
		Url:    c.Request().URL.Path,
		Method: c.Request().Method,
	}
}
