package utils

import (
	"fmt"
	"os"
)

// GetEnv returns the value of the environment variable named by the key. If the
// environment variable is not set or empty, an error is returned.
func GetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing environment variable \"%s\"", key)
	}
	return value, nil
}
