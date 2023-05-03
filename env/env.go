package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

// ErrEnvVarNotFound is an error that is raised when an environment variable is missing.
type ErrEnvVarNotFound string

func (envVar ErrEnvVarNotFound) Error() string {
	return fmt.Sprintf("%s was not found in the environment variables", string(envVar))
}

// ErrUnableToParseIntWithDefault is raises when converting a environment variable to int raises an error
type ErrUnableToParseIntWithDefault struct {
	key    string
	raw    string
	defVal int
}

func (e ErrUnableToParseIntWithDefault) Error() string {
	return fmt.Sprintf(
		"unable to parse .env variable '%s' with value '%s' as integer, setting to default '%d'",
		e.key,
		e.raw,
		e.defVal,
	)
}

// ErrUnableToParseInt is raises when converting a environment variable to int raises an error
type ErrUnableToParseInt struct {
	key string
	raw string
}

func (e ErrUnableToParseInt) Error() string {
	return fmt.Sprintf(
		"unable to parse .env variable '%s' with value '%s' as integer",
		e.key,
		e.raw,
	)
}

// ErrUnableToParseDurationWithDefault is raises when converting a environment variable to duration raises an error.
type ErrUnableToParseDurationWithDefault struct {
	key    string
	raw    string
	defVal time.Duration
}

func (e ErrUnableToParseDurationWithDefault) Error() string {
	return fmt.Sprintf(
		"unable to parse .env variable '%s' with value '%s' as duration, setting to default '%d'",
		e.key,
		e.raw,
		e.defVal,
	)
}

// ErrUnableToParseDuration is raises when converting a environment variable to int raises an error
type ErrUnableToParseDuration struct {
	key string
	raw string
}

func (e ErrUnableToParseDuration) Error() string {
	return fmt.Sprintf(
		"unable to parse .env variable '%s' with value '%s' as duration",
		e.key,
		e.raw,
	)
}

// InitEnv initializes the environment variables.
func InitEnv() error {
	paths := []string{}

	// Check to see if `.env` and `.env.default` exist before attempting to load environment vars from them.
	if _, err := os.Stat(".env"); err == nil {
		paths = append(paths, ".env")
	}
	if _, err := os.Stat(".env.default"); err == nil {
		paths = append(paths, ".env.default")
	}

	// Load the environment variables first from `.env.default` and then from `.env`, allowing `.env` to override when
	// necessary.
	err := godotenv.Load(paths...)
	if err != nil {
		return err
	}
	return nil
}

// InitEnvUnlessTest initializes the environment variables unless running in a test environment.
func InitEnvUnlessTest(envs ...string) error {
	if IsTest() {
		return nil
	}

	return InitEnv()
}

// Get simply returns the environment variable as a string, or an empty string when undefined.
func Get(key string) string {
	return os.Getenv(key)
}

// GetString returns the environment variable as a string, or the default value when undefined.
func GetString(key, defVal string) string {
	val := Get(key)
	if val == "" {
		return defVal
	}
	return val
}

// GetBool returns the environment variable as a bool, or the default value when undefined or unparsable.
func GetBool(key string, defVal bool) bool {
	val, err := strconv.ParseBool(Get(key))
	if err != nil {
		return defVal
	}
	return val
}

// GetBoolOrFalse returns the environment variable as a bool, or false when undefined or if it couldn't be parsed as a bool.
func GetBoolOrFalse(key string) bool {
	val, err := strconv.ParseBool(Get(key))
	if err != nil {
		return false
	}
	return val
}

// GetInt returns the environment variable as a int, or the default value when undefined.
func GetInt(key string, defVal int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return defVal
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return defVal
	}
	return val
}

// GetDuration returns the environment variable as a second based duration, or the default value when undefined.
func GetDuration(key string, defVal time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return defVal
	}
	duration, err := time.ParseDuration(raw)
	if err != nil {
		return defVal
	}
	return duration
}
