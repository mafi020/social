package env

import (
	"log"
	"os"
	"strconv"
)

func GetEnvOrPanic(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Panicf("Environment variable %s no found", key)
	}
	return val
}

func GetEnvAsIntOrPanic(key string) int {
	val := GetEnvOrPanic(key)
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Panicf("Invalid int value for %s: %v", key, err)
	}
	return intVal
}
