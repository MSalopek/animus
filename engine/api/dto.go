package api

import (
	"database/sql"
	"time"
)

type CreateKeyResponse struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	ClientKey    string `json:"client_key"`
	ClientSecret string `json:"client_secret"`
	Rights       string `json:"rights"`
	Disabled     bool   `json:"disabled"`

	CreatedAt time.Time    `json:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	ValidFrom time.Time    `json:"valid_from"`
	ValidTo   time.Time    `json:"valid_to"`
}

type UpdateKeyRequest struct {
	Rights   string `json:"rights"`
	Disabled bool   `json:"disabled"`
}
