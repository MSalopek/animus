package engine

import (
	"time"

	"github.com/msalopek/animus/model"
)

type CreateKeyResponse struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	ClientKey    string `json:"client_key"`
	ClientSecret string `json:"client_secret"`
	Rights       string `json:"rights"`
	Disabled     bool   `json:"disabled"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	ValidFrom time.Time  `json:"valid_from"`
}

type UpdateKeyRequest struct {
	Rights   *string `json:"rights,omitempty"`
	Disabled *bool   `json:"disabled,omitempty"`
}

type GetStorageResponse struct {
	Total    int                     `json:"total"`
	Returned int                     `json:"returned"`
	Rows     []*model.CountedStorage `json:"rows"`
}
