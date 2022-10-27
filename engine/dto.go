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
	Total    int              `json:"total"`
	Returned int              `json:"returned"`
	Rows     []*model.Storage `json:"rows"`
}

type UpdateUserRequest struct {
	Username  *string `json:"username"`
	Firstname *string `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Email     *string `json:"email"`

	WebhooksURL    *string `json:"webhooks_url"`
	WebhooksActive *bool   `json:"webhooks_active"`
}

type SyncAddFileResponse struct {
	Object   SyncAddFileObject `json:"object"`
	Status   string            `json:"status"`
	Error    string            `json:"error,omitempty"`
	RetryUrl string            `json:"retry_url,omitempty"`
}

type SyncAddFileObject struct {
	ID     int64   `json:"id" gorm:"primaryKey"`
	Cid    *string `json:"cid"`
	Dir    bool    `json:"dir"`
	Name   string  `json:"name"`
	Public bool    `json:"public"`
	Hash   *string `json:"hash,omitempty"`

	UploadStage *string `json:"stage"`
	Pinned      bool    `json:"pinned"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateNFTRequest struct {
	StorageID   int64  `json:"storage_id"`
	TokenID     int64  `json:"token_id"`
	ItemType    string `json:"item_type"`
	ExternalURL string `json:"external_url"`
}

type AfterMintWebhookRequest struct {
	ItemId      int64  `json:"itemId"`
	Type        string `json:"type"`
	ExternalURL string `json:"externalLink"`
	TokenID     int64  `json:"tokenId"`
}
