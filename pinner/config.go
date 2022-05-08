package pinner

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/msalopek/animus/storage"
)

type Config struct {
	NodeApiURL string `json:"node_api_url" yaml:"node_api_url"`
	LocalShell bool   `json:"local_shell" yaml:"local_shell"`

	DbDSN string `json:"db_dsn" yaml:"db_dsn"`
	DbURI string `json:"db_uri" yaml:"db_uri"`

	Debug    bool   `json:"debug" yaml:"debug"`
	TextLogs bool   `json:"text_logs" yaml:"text_logs"`
	LogFile  string `json:"log_file" yaml:"log_file"`

	Storage storage.Config `json:"storage" yaml:"storage"`
	Bucket  string         `json:"bucket" yaml:"bucket"`

	NsqLookupdURL  string   `json:"nsq_lookupd_url" yaml:"nsq_lookupd_url"`
	PublishTopics  []string `json:"publish_topic" yaml:"publish_topic"`
	SubscribeTopic string   `json:"subscribe_topic" yaml:"subscribe_topic"`

	MaxConcurrentRequests int `json:"max_concurrent_requests" yaml:"max_concurrent_requests"`
}

// Validate validates config.
func (c *Config) Validate() error {

	if len(c.DbDSN) == 0 && len(c.DbURI) == 0 {
		return errors.New("database connection info must be provided as URI or DSN")
	}

	if _, err := strconv.Atoi(strings.Split(c.NodeApiURL, ":")[1]); err != nil {
		return fmt.Errorf("error parsing node api url %s", err)
	}

	if _, err := strconv.Atoi(strings.Split(c.NsqLookupdURL, ":")[1]); err != nil {
		return fmt.Errorf("error parsing nsqd url %s", err)
	}

	return c.Storage.Validate()
}

type LoggingConfig struct {
	Pretty bool `yaml:"pretty"`
	Debug  bool `yaml:"debug"`
}
