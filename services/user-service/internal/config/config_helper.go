package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	trimmedValue := strings.TrimSpace(value)

	if !exists || trimmedValue == "" {
		return fallback
	}
	return trimmedValue
}

func getEnvRequired(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	trimmedValue := strings.TrimSpace(value)
	if !exists || trimmedValue == "" {
		return "", fmt.Errorf("%s is required ", key)
	}
	return trimmedValue, nil
}

func getEnvInt(key string, fallback int) (int, error) {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return fallback, nil
	}
	valueInt, err := strconv.Atoi(strings.TrimSpace(valueStr))
	if err != nil {
		return 0, fmt.Errorf("failed to convert %s to int", value)
	}
	return valueInt, nil
}

func getEnvBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	boolValue, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return boolValue
}

func validatePort(key, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%s must be an integer: %w", key, err)
	}
	if port < 0 || port > 65535 {
		return fmt.Errorf("%s must be between 0 and 65535", key)
	}
	return nil
}
