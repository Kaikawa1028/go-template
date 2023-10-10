package validate_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	CustomValidate "github.com/Kaikawa1028/go-template/app/validate"
	"github.com/stretchr/testify/assert"
)

/**
 400エラーの形式
 {
    "errors": [
        {
            "field": "companyName", //swaggerにて定義された変数名 (例：companyName など )
            "issue": "max",
            "option" : "3",
            "input" : "test"
        },
        ~~
    ],
    "url" : "v1/groups", //呼ばれたエンドポイント
    "method" : "Get"
}
*/
func TestBadRequest_400エラーの形式に沿ったレスポンスが返却が出来ること(t *testing.T) {
	validate := validator.New()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/sample", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	type TestBadRequest struct {
		Id   int    `validate:"required"`
		Name string `validate:"required"`
		Ids  []int  `validate:"unique"`
	}

	test := &TestBadRequest{
		Ids: []int{2, 2},
	}

	errs := validate.Struct(test)
	validateErrs := errs.(validator.ValidationErrors)

	badReq := CustomValidate.MakeBadRequestErrorResponse(c, validateErrs)

	assert.Len(t, badReq.Errors, 3)
	assert.Equal(t, badReq.Errors[0].Issue, "required")
	assert.Equal(t, badReq.Errors[0].Option, "")
	assert.Equal(t, badReq.Errors[0].Input, "0")
	assert.Equal(t, badReq.Errors[1].Issue, "required")
	assert.Equal(t, badReq.Errors[1].Option, "")
	assert.Equal(t, badReq.Errors[1].Input, "")
	assert.Equal(t, badReq.Errors[2].Issue, "unique")
	assert.Equal(t, badReq.Errors[2].Option, "")
	assert.Equal(t, badReq.Errors[2].Input, "2,2")
	assert.Equal(t, badReq.Method, "GET")
	assert.Equal(t, badReq.Url, "/sample")
}
