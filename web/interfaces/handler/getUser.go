package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/Kaikawa1028/go-template/app/domain/model"
)

type GetUser struct {
}

func NewGetUser() *GetUser {
	return &GetUser{}
}

type GetUserResponseUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`

	RoleOption GetUserResponseUserRoleOption `json:"roleOption"`
}

type GetUserResponseUserRoleOption struct {
	DisabledStoreEdit bool `json:"disabledStoreEdit"`
}

func (h GetUser) GetUser(c echo.Context) error {
	return c.JSON(http.StatusOK, h.buildResponse(c.Get("user").(*model.User)))
}

func (h GetUser) buildResponse(user *model.User) GetUserResponseUser {
	return GetUserResponseUser{
		Id:         user.ID,
		Username:   user.Username,
		Name:       user.Name,
		Email:      user.Email,
		Role:       user.Role,
		RoleOption: h.buildRoleOption(user),
	}
}

func (h GetUser) buildRoleOption(user *model.User) GetUserResponseUserRoleOption {
	roleOption := GetUserResponseUserRoleOption{
		DisabledStoreEdit: false,
	}

	if user.RoleOption != nil {
		roleOption.DisabledStoreEdit = user.RoleOption.DisabledStoreEdit
	}

	return roleOption
}
