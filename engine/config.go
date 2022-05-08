package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/msalopek/animus/storage"
)

type Config struct {
	AuthSecret          string `json:"auth_secret" yaml:"auth_secret"`
	AuthAuthority       string `json:"auth_authority" yaml:"auth_authority"`
	AuthExpirationHours int    `json:"auth_expiration_hours" yaml:"auth_expiration_hours"`
	HttpPort            string `json:"http_port" yaml:"http_port"`
	DbDSN               string `json:"db_dsn" yaml:"db_dsn"`
	DbURI               string `json:"db_uri" yaml:"db_uri"`
	Debug               bool   `json:"debug" yaml:"debug"`
	TextLogs            bool   `json:"text_logs" yaml:"text_logs"`
	LogFile             string `json:"log_file" yaml:"log_file"`

	NsqdURL  string `json:"nsqd_url" yaml:"nsqd_url"`
	NsqTopic string `json:"nsq_topic" yaml:"nsq_topic"`

	Bucket  string         `json:"bucket" yaml:"bucket"`
	Storage storage.Config `json:"storage" yaml:"storage"`
}

// Validate validates config.
func (c *Config) Validate() error {
	if utf8.RuneCountInString(c.AuthSecret) < 32 {
		return errors.New("signing secret must be at least 32 characters")
	}

	if c.AuthExpirationHours < 1 {
		return errors.New("token expiration hours must be at least 1 hour")
	}

	if len(c.AuthAuthority) < 1 {
		return errors.New("signing authority must be specified")
	}

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
