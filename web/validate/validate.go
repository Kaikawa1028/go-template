package validate

import (
	"regexp"
	"time"

	"github.com/Kaikawa1028/go-template/app/domain/system"

	"github.com/Kaikawa1028/go-template/app/const/date"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
	timer     system.ITimer
}

func NewCustomValidator(
	timer system.ITimer,
) *CustomValidator {
	v := validator.New()
	cv := CustomValidator{
		Validator: v,
		timer:     timer,
	}
	v.RegisterValidation("after", cv.after)
	v.RegisterValidation("after_or_equal", cv.afterEqual)
	v.RegisterValidation("before", cv.before)
	v.RegisterValidation("before_or_equal", cv.beforeEqual)
	v.RegisterValidation("regexp", cv.regexp)

	return &cv
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func (cv *CustomValidator) after(fl validator.FieldLevel) bool {
	field := fl.Field()
	fieldTime, err := time.Parse(date.FormatTypeDayTime, field.String())

	if err != nil {
		return false
	}

	param := fl.Param()
	targetTime, err := time.Parse(date.FormatTypeDayTime, param)

	if err != nil {
		return false
	}

	return fieldTime.After(targetTime)
}

func (cv *CustomValidator) afterEqual(fl validator.FieldLevel) bool {
	field := fl.Field()
	fieldTime, err := time.Parse(date.FormatTypeDayTime, field.String())

	if err != nil {
		return false
	}

	param := fl.Param()
	targetTime, err := time.Parse(date.FormatTypeDayTime, param)

	if err != nil {
		return false
	}

	return fieldTime.After(targetTime) || fieldTime.Equal(targetTime)
}

func (cv *CustomValidator) before(fl validator.FieldLevel) bool {
	field := fl.Field()
	fieldTime, err := time.Parse(date.FormatTypeDayTime, field.String())

	if err != nil {
		return false
	}

	param := fl.Param()
	targetTime, err := time.Parse(date.FormatTypeDayTime, param)

	if err != nil {
		return false
	}

	return fieldTime.Before(targetTime)
}

func (cv *CustomValidator) beforeEqual(fl validator.FieldLevel) bool {
	field := fl.Field()
	fieldTime, err := time.Parse(date.FormatTypeDayTime, field.String())

	if err != nil {
		return false
	}

	param := fl.Param()
	targetTime, err := time.Parse(date.FormatTypeDayTime, param)

	if err != nil {
		return false
	}

	return fieldTime.Before(targetTime) || fieldTime.Equal(targetTime)
}

func (cv *CustomValidator) regexp(fl validator.FieldLevel) bool {
	field := fl.Field()
	isMatch, err := regexp.MatchString(fl.Param(), field.String())

	if err != nil {
		return false
	}

	return isMatch
}
