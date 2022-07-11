package repo

import (
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const rawASCStorage = `WITH user_storage AS (
	SELECT *
	FROM  storage
	WHERE
	user_id = ?
	AND deleted_at IS NULL
)
SELECT *
FROM  (
	TABLE user_storage
	ORDER BY id ASC
	LIMIT ?
	OFFSET ?
) sub
RIGHT JOIN (
	SELECT count(*) FROM user_storage
) AS c (total_rows) ON true`

const rawDESCStorage = `WITH user_storage AS (
	SELECT *
	FROM  storage
	WHERE
	user_id = ?
	AND deleted_at IS NULL
)
SELECT *
FROM  (
	TABLE user_storage
	ORDER BY id DESC
	LIMIT ?
	OFFSET ?
) sub
RIGHT JOIN (
	SELECT count(*) FROM user_storage
) AS c (total_rows) ON true`

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

func (rpo *Repo) GetUserUploads(ctx QueryCtx, userID int) ([]*model.Storage, error) {
	var s []*model.Storage
	q := rpo.Where("user_id = ? AND deleted_at IS NULL", userID).
		Limit(ctx.Limit).
		Offset(ctx.Offset)
	if ctx.OrderBy != "" {
		q.Order(clause.OrderByColumn{
			Column: clause.Column{Name: ctx.OrderBy},
			Desc:   !ctx.Asc,
		})
	} else {
		q.Order(
			clause.OrderByColumn{
				Column: clause.Column{Name: "id"},
				Desc:   !ctx.Asc,
			})
	}
	res := q.Find(&s)
	if res.Error != nil {
		return nil, res.Error
	}
	return s, nil
}

// GetCountedUserUploads returns storage rows with total row count attached to each row.
// TODO: work out how to do keeping row counts in a smarter way.
func (rpo *Repo) GetCountedUserUploads(ctx QueryCtx, userID int) ([]*model.CountedStorage, error) {
	var s []*model.CountedStorage

	q := rpo.Raw(rawDESCStorage, userID, ctx.Limit, ctx.Offset)
	if ctx.Asc {
		q = rpo.Raw(rawASCStorage, userID, ctx.Limit, ctx.Offset)
	}

	res := q.Find(&s)
	if res.Error != nil {
		return nil, res.Error
	}
	return s, nil
}

func (rpo *Repo) GetUserUploadByCid(userID int, cid string) (*model.Storage, error) {
	var s model.Storage
	res := rpo.Where("user_id = ? AND cid = ? AND deleted_at IS NULL", userID, cid).Limit(1).Find(&s)
	if res.Error != nil {
		return nil, res.Error
	}
	return &s, nil
}

func (rpo *Repo) GetUserUploadByID(userID, id int) (*model.Storage, error) {
	var s model.Storage
	res := rpo.Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).First(&s)
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

func (rpo *Repo) DeleteUserUploadById(userID, id int) error {
	res := rpo.Where("user_id = ? AND id = ?", userID, id).Delete(&model.Storage{})

	return res.Error
}

func (rpo *Repo) GetApiClientByKey(key string) (*model.APIClient, error) {
	var c model.APIClient
	res := rpo.Table("users").
		Select("users.id AS user_id, users.email, keys.client_key, keys.client_secret").
		Joins("JOIN keys ON keys.user_id = users.id").
		Where(`keys.client_key = ?
			AND keys.deleted_at IS NULL
			AND keys.disabled IS false
			AND users.deleted_at IS NULL`, key).
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
	res := rpo.Where("id = ? AND user_id = ?", keyID, userID).Delete(&model.Key{})

	return res.Error
}

func (rpo *Repo) UpdateUserApiKey(userID int, keyID int, update *engine.UpdateKeyRequest) (*model.Key, error) {
	var key model.Key
	res := rpo.Table("keys").Where("id = ? AND user_id = ?", keyID, userID).First(&key)
	if res.Error != nil {
		return nil, res.Error
	}

	if update.Rights != nil {
		key.Rights = *update.Rights
	}

	if update.Disabled != nil {
		key.Disabled = *update.Disabled
	}

	res = rpo.Save(&key)
	if res.Error != nil {
		return nil, res.Error
	}
	return &key, nil
}
