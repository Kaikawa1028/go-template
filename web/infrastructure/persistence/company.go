package persistence

import (
	stdErrors "errors"

	"github.com/Kaikawa1028/go-template/app/domain/model"
	"github.com/Kaikawa1028/go-template/app/errors/types"

	"github.com/Kaikawa1028/go-template/app/domain/repository"
	"github.com/Kaikawa1028/go-template/app/errors"
	"gorm.io/gorm"
)

type Company struct {
}

func NewCompany() repository.Company {
	return &Company{}
}

func (c *Company) GetCompanyName(db *gorm.DB, id int) (company *model.Company, err error) {

	err = db.
		Select([]string{"name"}).
		Where("id = ?", id).
		First(&company).
		Error

	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			err = types.NewResourceNotFoundError()
		}
		return nil, errors.Wrap(err)
	}

	return company, nil
}
