package utilsconv

import "os"

func DefaultGetEnv(key string, default_ string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return default_
}
