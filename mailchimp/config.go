package mailchimp

import "errors"

var (
	ErrConfigEmptyDataCenter = errors.New("Config value 'DataCenter' was empty")
	ErrConfigEmptyAPIKey     = errors.New("Config value 'APIKey' was empty")
)

// Config represents the configuration for a MailChimp client instance.
type Config struct {
	DataCenter string
	APIKey     string
}

// Validate returns an error if there is any issue with the configuration, else
// returns nil if the configuration is considered "Valid".
func (c Config) Validate() error {
	if c.DataCenter == "" {
		return ErrConfigEmptyDataCenter
	}
	if c.APIKey == "" {
		return ErrConfigEmptyAPIKey
	}
	return nil
}
