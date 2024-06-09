package main

import (
	"fmt"
	"strconv"
	"strings"
)

func boolSetting(settings map[string]string, key string) (result bool, err error) {
	value := settings[key]
	switch strings.ToLower(value) {
	case "true":
		return true, nil

	case "false", "":
		return false, nil

	default:
		return false, fmt.Errorf("invalid boolean setting %s: %s", key, value)
	}
}

func intSetting(settings map[string]string, key string, defaultValue int) (result int, err error) {
	value := settings[key]
	if value == "" {
		return defaultValue, nil
	}
	result, err = strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer setting %s: %s", key, value)
	}
	return
}

func requireSetting(settings map[string]string, key string) (result string, err error) {
	result = settings[key]
	if result == "" {
		err = fmt.Errorf("missing required setting %s", key)
	}
	return
}
