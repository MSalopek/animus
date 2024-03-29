package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

const (
	UploadStageIPFS    = "ipfs"
	UploadStageStorage = "storage"
)

type ClientAccess string

const (
	ClientAccessRead            = "r"
	ClientAccessReadWrite       = "rw"
	ClientAccessReadWriteDelete = "rwd"

	TokenTypeRegisterEmail = "register_email"
	TokenTypeResetPass     = "reset_pass"

	DefaultMaxKeys = 5
)

type User struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	Username  string     `json:"username"`
	Firstname *string    `json:"firstname"`
	Lastname  *string    `json:"lastname"`
	Email     string     `json:"email"`
	Password  []byte     `json:"-"`
	MaxKeys   int32      `json:"max_keys"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Verified  bool       `json:"verified"`

	WebhooksURL    *string `json:"webhooks_url,omitempty"`
	WebhooksActive *bool   `json:"webhooks_active"`
}

type Storage struct {
	ID       int64          `json:"id" gorm:"primaryKey"`
	Cid      *string        `json:"cid"`
	Dir      bool           `json:"dir"`
	UserID   int64          `json:"-"`
	Name     string         `json:"name"`
	Public   bool           `json:"public"`
	Metadata datatypes.JSON `json:"meta"`
	Hash     *string        `json:"hash,omitempty"`

	UploadStage   *string `json:"stage"`
	StorageBucket *string `json:"-"`
	StorageKey    *string `json:"-"`
	Pinned        bool    `json:"pinned"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// CountedStorage is Storage with count column added.
// count is usefull in pagination; it is faster to return
// the count than it is to make 2 separate queries to fetch and count data.
type CountedStorage struct {
	Storage
	TotalRows int `json:"-"`
}

// ToStorage returns the underlying Storage record from CountedStorage.
func (cs *CountedStorage) ToStorage() *Storage {
	return &cs.Storage
}

func (Storage) TableName() string {
	return "storage"
}

type Gateway struct {
	ID       int64         `json:"id" gorm:"primaryKey"`
	UserID   sql.NullInt64 `json:"user_id"`
	Name     string        `json:"name"`
	Slug     string        `json:"slug"`
	PublicID uuid.UUID     `json:"public_id"`
}

type Key struct {
	ID           int64  `json:"id" gorm:"primaryKey"`
	UserID       int64  `json:"user_id"`
	ClientKey    string `json:"client_key"`
	ClientSecret string `json:"-"`
	Rights       string `json:"rights"`
	Disabled     bool   `json:"disabled"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	ValidFrom time.Time  `json:"valid_from"`
	ValidTo   PgTime     `json:"-"`
}

type Token struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	UserID int64  `json:"user_id"`
	Type   string `json:"type"`
	Token  string `json:"token"`

	ValidFrom time.Time `json:"valid_from"`
	ValidTo   PgTime    `json:"-"`
	IsUsed    bool      `json:"is_used"`
}

type Subscription struct {
	ID        int64  `json:"-" gorm:"primaryKey"`
	PublicID  string `json:"public_id"`
	Name      string `json:"name"`
	Promotion bool   `json:"promotion"`

	Price    decimal.Decimal `json:"price"`
	Currency string          `json:"currency"`

	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   PgTime    `json:"-"`
}

type UserSubscription struct {
	ID             int64  `json:"-" gorm:"primaryKey"`
	PublicID       string `json:"public_id"`
	Promotion      bool   `json:"promotion"`
	SubscriptionID int64  `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   PgTime    `json:"-"`
}

type APIClient struct {
	UserID       int64
	Email        string
	ClientKey    string
	ClientSecret string
}
