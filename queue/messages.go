package queue

import "encoding/json"

// Holds path to storage object and DB StorageID.
type PinRequest struct {
	StorageID int    `json:"storage_id,omitempty"` // db storage pk
	Dir       bool   `json:"dir,omitempty"`
	Key       string `json:"key,omitempty"` // path in object storage (s3, minio or local path)
}

func (pr *PinRequest) Unmarshal(raw []byte) error {
	return json.Unmarshal(raw, pr)
}

func (pr *PinRequest) Marshal() ([]byte, error) {
	return json.Marshal(pr)
}
