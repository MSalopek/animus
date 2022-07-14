package client

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/msalopek/animus/storage"
)

type Config struct {
	HttpPort string `json:"http_port,omitempty" yaml:"http_port"`
	DbDSN    string `json:"db_dsn,omitempty" yaml:"db_dsn"`
	DbURI    string `json:"db_uri,omitempty" yaml:"db_uri"`
	Debug    bool   `json:"debug,omitempty" yaml:"debug"`
	TextLogs bool   `json:"text_logs,omitempty" yaml:"text_logs"`
	LogFile  string `json:"log_file,omitempty" yaml:"log_file"`

	NsqdURL        string `json:"nsqd_url,omitempty" yaml:"nsqd_url"`
	NsqPinnerTopic string `json:"nsq_pinner_topic,omitempty" yaml:"nsq_pinner_topic"`

	Bucket  string         `json:"bucket,omitempty" yaml:"bucket"`
	Storage storage.Config `json:"storage,omitempty" yaml:"storage"`
	GinMode string         `json:"gin_mode,omitempty" yaml:"gin_mode"`
}

// Validate validates config.
func (c *Config) Validate() error {

	if len(c.DbDSN) == 0 && len(c.DbURI) == 0 {
		return errors.New("database connection info must be provided as URI or DSN")
	}

	if len(c.Bucket) == 0 {
		return errors.New("bucket name must be configured")
	}

	if _, err := strconv.Atoi(strings.Split(c.HttpPort, ":")[1]); err != nil {
		return fmt.Errorf("error parsing http port %s", err)
	}

	if _, err := strconv.Atoi(strings.Split(c.NsqdURL, ":")[1]); err != nil {
		return fmt.Errorf("error parsing nsqd url %s", err)
	}

	return c.Storage.Validate()
}

type LoggingConfig struct {
	Pretty bool `yaml:"pretty"`
	Debug  bool `yaml:"debug"`
}
