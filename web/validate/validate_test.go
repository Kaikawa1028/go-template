package validate_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	CutomValidate "github.com/Kaikawa1028/go-template/app/validate"
	"github.com/stretchr/testify/assert"
)

func TestAfter_フィールドの値が与えられた日付より後であるかバリデーションするテスト(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	validate := CutomValidate.NewCustomValidator(cont.MockTimer)

	type TestAfter struct {
		Date string `validate:"after=2021-01-31 23:59:59"`
	}

	tSuccess := &TestAfter{
		Date: "2022-01-31 23:59:59",
	}

	errs := validate.Validate(tSuccess)

	assert.Equal(t, errs, nil)

	tFail := &TestAfter{
		Date: "2021-01-31 23:59:59",
	}

	errs = validate.Validate(tFail)
	validateErrs := errs.(validator.ValidationErrors)

	assert.Len(t, validateErrs, 1)
	assert.Equal(t, validateErrs[0].Tag(), "after")
}

func TestAfterEqual_フィールドが指定した日付以降であることをバリデートするテスト(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	validate := CutomValidate.NewCustomValidator(cont.MockTimer)

	type TestAfter struct {
		Date string `validate:"after_or_equal=2021-01-31 23:59:59"`
	}

	tSuccess := &TestAfter{
		Date: "2021-01-31 23:59:59",
	}

	errs := validate.Validate(tSuccess)

	assert.Equal(t, errs, nil)

	tFail := &TestAfter{
		Date: "2021-01-30 23:59:59",
	}

	errs = validate.Validate(tFail)
	validateErrs := errs.(validator.ValidationErrors)

	assert.Len(t, validateErrs, 1)
	assert.Equal(t, validateErrs[0].Tag(), "after_or_equal")
}

func TestBefore_フィールドが指定された日付より前であることをバリデートするテスト(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	validate := CutomValidate.NewCustomValidator(cont.MockTimer)

	type TestBefore struct {
		Date string `validate:"before=2021-01-31 23:59:59"`
	}

	tSuccess := &TestBefore{
		Date: "2021-01-30 23:59:59",
	}

	errs := validate.Validate(tSuccess)

	assert.Equal(t, errs, nil)

	tFail := &TestBefore{
		Date: "2021-01-31 23:59:59",
	}

	errs = validate.Validate(tFail)
	validateErrs := errs.(validator.ValidationErrors)

	assert.Len(t, validateErrs, 1)
	assert.Equal(t, validateErrs[0].Tag(), "before")
}

func TestBeforeEqual_フィールドが指定した日付以降であることをバリデートするテスト(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	validate := CutomValidate.NewCustomValidator(cont.MockTimer)

	type TestBefore struct {
		Date string `validate:"before_or_equal=2021-01-31 23:59:59"`
	}

	tSuccess := &TestBefore{
		Date: "2021-01-31 23:59:59",
	}

	errs := validate.Validate(tSuccess)

	assert.Equal(t, errs, nil)

	tFail := &TestBefore{
		Date: "2021-02-01 23:59:59",
	}

	errs = validate.Validate(tFail)
	validateErrs := errs.(validator.ValidationErrors)

	assert.Len(t, validateErrs, 1)
	assert.Equal(t, validateErrs[0].Tag(), "before_or_equal")
}

func TestRegexp_正規表現に合致するかバリデートするテスト(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	validate := CutomValidate.NewCustomValidator(cont.MockTimer)

	type TestRegexp struct {
		CustomID string `validate:"regexp=^test"`
	}

	tSuccess := &TestRegexp{
		CustomID: "testId",
	}

	errs := validate.Validate(tSuccess)

	assert.Equal(t, errs, nil)

	tFail := &TestRegexp{
		CustomID: "sampleId",
	}

	errs = validate.Validate(tFail)
	validateErrs := errs.(validator.ValidationErrors)

	assert.Len(t, validateErrs, 1)
	assert.Equal(t, validateErrs[0].Tag(), "regexp")
}

func TestRegexp_正規表現に合致するかバリデートするテスト2(t *testing.T) {
	cont := beforeEach(t)
	defer afterEach(cont)
	validate := CutomValidate.NewCustomValidator(cont.MockTimer)

	type TestRegexp struct {
		CustomID string `validate:"regexp=^[0-9a-zA-Z-_]+$"`
	}

	tSuccess := &TestRegexp{
		CustomID: "testI-_d",
	}

	errs := validate.Validate(tSuccess)

	assert.Equal(t, errs, nil)

	tFail := &TestRegexp{
		CustomID: "sampleI*/..d",
	}

	errs = validate.Validate(tFail)
	validateErrs := errs.(validator.ValidationErrors)

	assert.Len(t, validateErrs, 1)
	assert.Equal(t, validateErrs[0].Tag(), "regexp")
}
