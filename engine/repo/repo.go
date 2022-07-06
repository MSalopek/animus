package repo

import (
	"time"

	"github.com/msalopek/animus/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	*gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{DB: db}
}

func (rpo *Repo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	res := rpo.Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (rpo *Repo) GetUserUploads(ctx QueryCtx, userID int) ([]model.Storage, error) {
	var s []model.Storage
	q := rpo.Where("user_id = ?", userID).
		Limit(ctx.Limit).
		Offset(ctx.Offset)
	if ctx.OrderBy != "" {
		q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx.OrderBy},
			Desc:   !ctx.Asc},
		)

	}
	res := q.Find(&s)
	if res.Error != nil {
		return nil, res.Error
	}
	return s, nil
}

func (rpo *Repo) GetUserUploadByCid(userID int, cid string) (*model.Storage, error) {
	var s model.Storage
	res := rpo.Where("user_id = ? AND cid = ?", userID, cid).Limit(1).Find(&s)
	if res.Error != nil {
		return nil, res.Error
	}
	return &s, nil
}

func (rpo *Repo) CreateStorage(s *model.Storage) error {
	res := rpo.Create(s)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (rpo *Repo) GetApiClientByKey(key string) (*model.APIClient, error) {
	var c model.APIClient
	res := rpo.Table("users").
		Select("users.id, users.email, keys.client_key, keys.client_secret").
		Joins("JOIN keys ON keys.user_id = user.id").
		Where("keys.client_key = ? AND key.deleted_at IS NULL", key).
		Scan(&c)

	if res.Error != nil {
		return nil, res.Error
	}
	return &c, nil
}

func (rpo *Repo) GetUserApiKeys(ctx QueryCtx, userID int) ([]*model.Key, error) {
	var keys []*model.Key
	q := rpo.Where("user_id = ? AND deleted_at IS NULL", userID).
		Limit(ctx.Limit).
		Offset(ctx.Offset)
	if ctx.OrderBy != "" {
		q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx.OrderBy},
			Desc:   !ctx.Asc},
		)

	}
	res := q.Find(&keys)
	if res.Error != nil {
		return nil, res.Error
	}
	return keys, nil
}

func (rpo *Repo) DeleteUserApiKey(userID int, keyID int) error {
	res := rpo.Table("keys").Where("id = ? AND user_id = ?", keyID, userID).
		Update("deleted_at", time.Now())

	return res.Error
}

func (rpo *Repo) UpdateUserApiKey(userID int, keyID int, upd *model.Key) error {
	res := rpo.Table("keys").Where("id = ? AND user_id = ?", keyID, userID).
		Updates(*upd)

	return res.Error
}
