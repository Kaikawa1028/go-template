package model

import (
	"github.com/Kaikawa1028/go-template/app/const/app"
	"time"
)

type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Role      string    `db:"slug"`
	Token     string    `db:"api_token"`
	ExpiresAt time.Time `db:"expires_at"`

	// RoleOption レコードが存在しない場合はnilが格納されます
	RoleOption *AdminUserRoleOption
}

func (u User) IsAdministrator() bool {
	return u.Role == app.RoleAdministrator
}

func (u User) IsAgency() bool {
	return u.Role == app.RoleAgency
}

func (u User) IsEditor() bool {
	return u.Role == app.RoleEditor
}

func (u User) IsViewer() bool {
	return u.Role == app.RoleViewer
}

func (u User) IsManager() bool {
	return u.Role == app.RoleManager
}

func (u User) IsWorker() bool {
	return u.Role == app.RoleWorker
}

func (u User) IsWorkerCanEditStore() bool {
	if u.Role != app.RoleWorker {
		return false
	}

	if u.RoleOption == nil {
		return true
	}

	return !u.RoleOption.DisabledStoreEdit
}
