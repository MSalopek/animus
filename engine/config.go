package engine

import (
	"errors"
	"unicode/utf8"
)

type Config struct {
	AuthSecret       string        `yaml:"signing_secret"`     // JTW  claims signing key
	AuthAuthority    string        `yaml:"signing_authority"`  // JTW claims signing authority
	DatabaseDSN      string        `yaml:"db_dsn"`             // database connection DSN
	HTTPPort         string        `yaml:"http_port"`          // Animus engine HTTP port (exposed to public)
	IPFSRPCURL       string        `yaml:"ipfs_rpc_url"`       // IPFS node private HTTP-RPC-API URL (not exposed to public)
	LocalStoragePath string        `yaml:"local_storage_path"` // path to local file storage on disk (cache directory)
	Logging          LoggingConfig `yaml:"logging"`            // logging configuration -> print JSON vs formatted text and allow debug mode
}

// Validate validates config.
func (c *Config) Validate() error {
	if utf8.RuneCountInString(c.AuthSecret) < 32 {
		return errors.New("signing secret must be at least 32 characters")
	}

	if len(c.AuthAuthority) < 1 {
		return errors.New("signing authority must be specified")
	}
	return nil
}

type LoggingConfig struct {
	Pretty bool `yaml:"pretty"`
	Debug  bool `yaml:"debug"`
}
