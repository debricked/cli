package util

import "os"

func EnvKeyIsSet(key string) bool {
	value, isPresent := os.LookupEnv(key)
	if isPresent && len(value) > 0 {
		return true
	}
	return false
}
