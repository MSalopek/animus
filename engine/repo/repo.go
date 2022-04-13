package repo

import (
	"github.com/msalopek/animus/model"
	"gorm.io/gorm"
)

type Repo struct {
	*gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{DB: db}
}

func (rpo *Repo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	res := rpo.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (rpo *Repo) CreateStorage(s *model.Storage) error {
	res := rpo.Create(s)
	if res.Error != nil {
		return res.Error
	}
	return nil
}