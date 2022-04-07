package dbclient

import validation "github.com/go-ozzo/ozzo-validation"

type Config struct {
	ServiceName string
	Host        string
	User        string
	Name        string
	Password    string
	MaxConns    int
	Port        int
	SearchPath  string
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ServiceName, validation.Required, validation.Required.Error("SERVICE_NAME was not specified in env")),
		validation.Field(&c.Name, validation.Required, validation.Required.Error("DB_NAME was not specified in env")),
	)
}
