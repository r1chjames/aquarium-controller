package utils

import (
	"fmt"
	"os"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		if len(fallback) > 0 {
			return fallback
		} else {
			panic(fmt.Sprintf("Required application variable %s not defined", key))
		}
	}
	return value
}
