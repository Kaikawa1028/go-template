package usecase

import (
	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"github.com/Kaikawa1028/go-template/app/errors"
	"gorm.io/gorm"
)

type GetCompanyName struct {
	db       *gorm.DB
	userRepo repository.User
	company  repository.Company
}

func NewGetCompanyName(db *gorm.DB, userRepo repository.User, company repository.Company) *GetCompanyName {
	return &GetCompanyName{
		db:       db,
		userRepo: userRepo,
		company:  company,
	}
}

func (u GetCompanyName) GetCompanyName(user *model.User, companyId int) (*model.Company, error) {
	company, err := u.company.GetCompanyName(u.db, companyId)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return company, nil
}
