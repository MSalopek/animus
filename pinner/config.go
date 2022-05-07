package pinner

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/msalopek/animus/storage"
)

type Config struct {
	NodeApiURL string `json:"node_api_url,omitempty" yaml:"node_api_url"`
	LocalShell bool   `json:"local_shell,omitempty" yaml:"local_shell,omitempty"`

	DbDSN string `json:"db_dsn,omitempty" yaml:"db_dsn,omitempty"`
	DbURI string `json:"db_uri,omitempty" yaml:"db_uri,omitempty"`

	Debug    bool   `json:"debug,omitempty" yaml:"debug,omitempty"`
	TextLogs bool   `json:"text_logs,omitempty" yaml:"text_logs,omitempty"`
	LogFile  string `json:"log_file,omitempty" yaml:"log_file,omitempty"`

	Storage storage.Config `json:"storage,omitempty" yaml:"storage"`
	Bucket  string         `json:"bucket,omitempty" yaml:"bucket"`

	NsqLookupdURL  string   `json:"nsq_lookupd_url,omitempty" yaml:"nsq_lookupd_url"`
	PublishTopics  []string `json:"publish_topic,omitempty" yaml:"publish_topic"`
	SubscribeTopic string   `json:"subscribe_topic,omitempty" yaml:"subscribe_topic"`

	// defines buffer size for request chan
	QueueSize int `json:"queue_depth" yaml:"queue_depth"`
}

// Validate validates config.
func (c *Config) Validate() error {

	if len(c.DbDSN) == 0 && len(c.DbURI) == 0 {
		return errors.New("database connection info must be provided as URI or DSN")
	}

	if _, err := strconv.Atoi(strings.Split(c.NodeApiURL, ":")[1]); err != nil {
		return errors.New(fmt.Sprintf("error parsing node api url %s", err))
	}

	return c.Storage.Validate()
}

type LoggingConfig struct {
	Pretty bool `yaml:"pretty"`
	Debug  bool `yaml:"debug"`
}
