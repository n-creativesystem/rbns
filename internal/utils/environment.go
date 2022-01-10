package utils

import (
	"os"
	"strings"
	"sync"
)

var (
	keys []string
	once sync.Once
)

func GetEnvKeys() []string {
	once.Do(func() {
		keys = make([]string, 0, 100)
		for _, e := range os.Environ() {
			kv := strings.SplitN(e, "=", 2)
			keys = append(keys, kv[0])
		}
	})
	return keys
}
