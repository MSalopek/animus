package queue

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
)

const (
	SourceUser   = "USER"
	SourceClient = "CLIENT"

	MailerTypeRegister  MailerMessageType = "register"
	MailerTypeResetPass MailerMessageType = "reset"
)

var validMailerMessage = map[MailerMessageType]struct{}{
	MailerTypeRegister:  struct{}{},
	MailerTypeResetPass: struct{}{},
}

// Holds path to storage object and DB StorageID.
type PinRequest struct {
	UserID    int    `json:"user_id,omitempty"`
	StorageID int    `json:"storage_id,omitempty"` // db storage pk
	CID       string `json:"cid,omitempty"`        // IPFS cid
	Dir       bool   `json:"dir,omitempty"`
	Key       string `json:"key,omitempty"`   // path in object storage (s3, minio or local path)
	Source    string `json:"src,omitempty"`   // metadata about source that sent the PinRequest
	Unpin     bool   `json:"unpin,omitempty"` // if true the CID will be unpinned
}

func (pr *PinRequest) Unmarshal(raw []byte) error {
	return json.Unmarshal(raw, pr)
}

func (pr *PinRequest) Marshal() ([]byte, error) {
	return json.Marshal(pr)
}

type MailerMessageType string

type MailerMessage struct {
	Type      MailerMessageType `json:"type"`
	Username  *string           `json:"username,omitempty"`
	Firstname *string           `json:"firstname,omitempty"`
	Lastname  *string           `json:"lastname,omitempty"`
	URL       string            `json:"url"`
	Email     string            `json:"email"`
}

func (m *MailerMessage) Validate() error {
	if _, ok := validMailerMessage[m.Type]; !ok {
		return fmt.Errorf("invalid mailer message type: %s", m.Type)
	}

	if m.Type == MailerTypeRegister && (m.Email == "" || m.URL == "") {
		return fmt.Errorf("registration messages must include email and url")
	}

	if m.Type == MailerTypeResetPass && (m.Email == "" || m.URL == "") {
		return fmt.Errorf("password reset messages must include email and url")
	}
	return nil
}

func (m *MailerMessage) Unmarshal(raw []byte) error {
	return json.Unmarshal(raw, m)
}

func (m *MailerMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

type WebhookMessage struct {
	StorageID   int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	Cid         *string        `json:"cid"`
	Name        string         `json:"name"`
	Metadata    datatypes.JSON `json:"meta"`
	UploadStage *string        `json:"stage"`
	Pinned      bool           `json:"pinned"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
	Status      string         `json:"status"`
}

func (w *WebhookMessage) Unmarshal(raw []byte) error {
	return json.Unmarshal(raw, w)
}

func (w *WebhookMessage) Marshal() ([]byte, error) {
	return json.Marshal(w)
}
