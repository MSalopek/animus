package storage

import (
	"errors"
)

// TODO: write some integration tests for this

type Config struct {
	Lifetime     int    `json:"lifetime" yaml:"lifetime"` // DAYS before bucket expiration
	Region       string `json:"region" yaml:"region"`
	Endpoint     string `json:"endpoint" yaml:"endpoint"`
	TmpDir       string `json:"tmp_dir,omitempty" yaml:"tmp_dir,omitempty"`               // path to temp directory on disk for storing downloads
	TmpDirPrefix string `json:"tmp_dir_prefix,omitempty" yaml:"tmp_dir_prefix,omitempty"` // prefix for files and dirs in temp directory

	AccessKeyID     string `json:"access_key_id" yaml:"access_key_id"`         // admin/root access credentials
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret"` // admin/root access credentials
	AccessToken     string `json:"access_token,omitempty" yaml:"access_token,omitempty"`
}

func (c Config) Validate() error {
	if c.Lifetime == 0 {
		return errors.New("invalid storage Lifetime")
	}
	if len(c.Region) == 0 {
		return errors.New("invalid storage Region")
	}
	if len(c.Endpoint) == 0 {
		return errors.New("invalid storage Endpoint")
	}
	if len(c.AccessKeyID) == 0 {
		return errors.New("invalid storage AccessKeyID")
	}
	if len(c.AccessKeySecret) == 0 {
		return errors.New("invalid storage AccessKeySecret")
	}
	return nil
}
