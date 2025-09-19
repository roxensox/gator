package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Global config filename variable
const configFileName = ".gatorconfig.json"

// Global permission code
const permCode = 0o600

func getConfigFilePath() (string, error) {
	// Returns the path to the config file on the user's system

	// Gets the home directory path
	homeDir, err := os.UserHomeDir()
	// Returns any error
	if err != nil {
		return "", err
	}

	// Constructs and returns the path to the config file
	return fmt.Sprintf("%s/%s", homeDir, configFileName), nil
}

func write(cfg Config) error {
	// Writes the Config object to a JSON file

	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, permCode)
	return err
}
