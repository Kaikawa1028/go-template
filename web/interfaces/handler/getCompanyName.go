package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/errors"
	"github.com/Kaikawa1028/go-template/app/usecase"
)

type GetCompanyName struct {
	useCase *usecase.GetCompanyName
}

func NewGetCompanyName(useCase *usecase.GetCompanyName) *GetCompanyName {
	return &GetCompanyName{
		useCase: useCase,
	}
}

type GetCompanyNameRequest struct {
	CompanyId int `param:"company-id"`
}

type GetCompanyNameResponse struct {
	Name string `json:"name"`
}

func (h GetCompanyName) GetCompanyName(c echo.Context) error {
	var r GetCompanyNameRequest

	if err := c.Bind(&r); err != nil {
		return errors.Wrap(err)
	}

	user := c.Get("user").(*model.User)

	company, err := h.useCase.GetCompanyName(user, r.CompanyId)

	if err != nil {
		return errors.Wrap(err)
	}

	response := h.buildResponse(company)
	return c.JSON(http.StatusOK, response)
}

func (h GetCompanyName) buildResponse(company *model.Company) GetCompanyNameResponse {

	return GetCompanyNameResponse{
		Name: company.Name,
	}
}
