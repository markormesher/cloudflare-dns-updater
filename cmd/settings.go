package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func getSettings() ([]ZoneSettings, error) {
	filePath := os.Getenv("SETTINGS_FILE")
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error loading settings: %w", err)
	}

	var settings []ZoneSettings
	err = json.Unmarshal(raw, &settings)
	if err != nil {
		return nil, fmt.Errorf("error loading settings: %w", err)
	}

	return settings, nil
}

func getCheckInterval() (int, error) {
	raw := os.Getenv("CHECK_INTERVAL_SECONDS")
	if raw == "" {
		return 0, nil
	}

	interval, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid check interval: %w", err)
	}
	if interval < 0 {
		return 0, fmt.Errorf("invalid check interval: %v", interval)
	}

	return interval, nil
}
