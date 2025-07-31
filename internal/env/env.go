package env

import (
	"log"
	"os"
)

func GetEnvOrPanic(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Panicf("Environment variable %s no found", key)
	}
	return val
}