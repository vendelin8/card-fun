package util

import (
	"os"
)

// GetEnv returns the value of an environment variable by the given key, or the
// given default vale, if it doesn't exist.
func GetEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
