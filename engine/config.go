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
	AuthSecret          string `json:"auth_secret,omitempty" yaml:"auth_secret"`
	AuthAuthority       string `json:"auth_authority,omitempty" yaml:"auth_authority"`
	AuthExpirationHours int    `json:"auth_expiration_hours,omitempty" yaml:"auth_expiration_hours"`
	HttpPort            string `json:"http_port,omitempty" yaml:"http_port"`
	DbDSN               string `json:"db_dsn,omitempty" yaml:"db_dsn"`
	DbURI               string `json:"db_uri,omitempty" yaml:"db_uri,omitempty"`
	Debug               bool   `json:"debug,omitempty" yaml:"debug,omitempty"`
	TextLogs            bool   `json:"text_logs,omitempty" yaml:"text_logs,omitempty"`
	LogFile             string `json:"log_file,omitempty" yaml:"log_file,omitempty"`

	Bucket  string         `json:"bucket,omitempty" yaml:"bucket"`
	Storage storage.Config `json:"storage,omitempty" yaml:"storage"`
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
		return errors.New(fmt.Sprintf("error parsing http port %s", err))
	}

	return c.Storage.Validate()
}

type LoggingConfig struct {
	Pretty bool `yaml:"pretty"`
	Debug  bool `yaml:"debug"`
}
