package queue

import "encoding/json"

const (
	SourceUser   = "USER"
	SourceClient = "CLIENT"
)

// Holds path to storage object and DB StorageID.
type PinRequest struct {
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
