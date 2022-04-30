package queue

// Holds path storage object and DB StorageID.
type StorageMessage struct {
	Topic     string `json:"topic,omitempty"`      // nsq topic
	StorageID int    `json:"storage_id,omitempty"` // db storage pk
	Dir       bool   `json:"dir,omitempty"`
	Path      string `json:"path,omitempty"` // path in object storage (s3, minio or local path)
}
