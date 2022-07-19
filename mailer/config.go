package mailer

import (
	"errors"
)

type Config struct {
	// mailer REST API key
	ApiKey string `json:"api_key,omitempty" yaml:"api_key"`
	// mailer REST API secret
	ApiSecret string `json:"api_secret,omitempty" yaml:"api_secret"`
	// NSQ email topic
	Topic string `json:"topic,omitempty" yaml:"topic"`
	// nsq lookup for nsq connections
	NsqLookup string `json:"nsq_lookup,omitempty" yaml:"nsq_lookup"`

	// default sender
	Sender     string `json:"sender,omitempty" yaml:"sender"`
	SenderName string `json:"sender_name,omitempty" yaml:"sender_name"`

	// output emails to stdout (does not send API requests)
	DryRun bool `json:"dry_run,omitempty" yaml:"dry_run"`
	Debug  bool `json:"debug,omitempty" yaml:"debug"`
}

// Validate validates config.
func (c *Config) Validate() error {

	if c.ApiKey == "" {
		return errors.New("api key required")
	}
	if c.ApiSecret == "" {
		return errors.New("api secret required")
	}
	if c.Topic == "" {
		return errors.New("topic required")
	}
	if c.NsqLookup == "" {
		return errors.New("nsq lookupd address required")
	}

	return nil
}
