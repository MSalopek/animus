package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	Firstname *string    `json:"firstname"`
	Lastname  *string    `json:"lastname"`
	Email     string     `json:"email"`
	Password  []byte     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Storage struct {
	ID        int64      `json:"id"`
	Cid       *string    `json:"cid"`
	UserID    int64      `json:"user_id"`
	Name      string     `json:"name"`
	Public    bool       `json:"public"`
	Metadata  string     `json:"metadata"`
	Local     bool       `json:"local"`
	LocalPath *string    `json:"local_path"`
	Hash      *string    `json:"hash"`
	Pinned    bool       `json:"pinned"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Storage) TableName() string {
	return "storage"
}

type Gateway struct {
	ID       int64         `json:"id"`
	UserID   sql.NullInt64 `json:"user_id"`
	Name     string        `json:"name"`
	Slug     string        `json:"slug"`
	PublicID uuid.UUID     `json:"public_id"`
}
