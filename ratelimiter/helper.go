package ratelimiter

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var StFormat = "2006-01-02 15:04:05"

func PrintfD(config *LimiterConfig, format string, keyType string, key string, a ...any) (n int, err error) {
	if config.Debug {
		timeString := time.Now().UTC().Format(StFormat)
		args := []any{timeString, keyType, key}
		args = append(args, a...)
		return fmt.Printf("%s [RATE LIMITER][%s][%s] "+format+"\n", args...)
	}

	return 0, nil
}

func PrintfWD(config *LimiterConfig, format string, a ...any) (n int, err error) {
	if config.Debug {
		timeString := time.Now().UTC().Format(StFormat)
		args := []any{timeString}
		args = append(args, a...)
		return fmt.Printf("%s [RATE LIMITER] "+format+"\n", args...)
	}

	return 0, nil
}

func GetBlockTime(block *time.Time) float64 {
	return time.Until(*block).Seconds()
}

func GetEnvString(key string) (string, bool) {
	value, encontrou := os.LookupEnv(key)
	if !encontrou {
		return "", false
	}
	if value == "" {
		return "", false
	}
	return value, true
}

func GetEnvBoolean(key string) (bool, bool) {
	value, encontrou := os.LookupEnv(key)
	if !encontrou {
		return false, false
	}
	if value == "" {
		return false, false
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, false
	}
	return parsed, true
}

func GetEnvLargeint(key string) (int64, bool) {
	value, encontrou := os.LookupEnv(key)
	if !encontrou {
		return 0, false
	}
	if value == "" {
		return 0, false
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}
