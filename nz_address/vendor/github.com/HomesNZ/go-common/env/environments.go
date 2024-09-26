package env

import "os"

// Env returns the environment the application is running in.
func Env() string {
	return os.Getenv("ENV")
}

// IsDev returns true when the environment is `dev`
func IsDev() bool {
	return Env() == "development"
}

// IsTest returns true when the environment is `test`
func IsTest() bool {
	return Env() == "test"
}

// IsProd returns true when the environment is `prod`
func IsProd() bool {
	return Env() == "production"
}
