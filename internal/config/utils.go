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

	// Gets the path to the config file
	configPath, err := getConfigFilePath()
	// Returns any errors
	if err != nil {
		return err
	}

	// Marshals the input Config object to json data
	data, err := json.Marshal(cfg)
	// Returns any errors
	if err != nil {
		return err
	}

	// Writes the file at the configPath location
	err = os.WriteFile(configPath, data, permCode)
	// Returns any errors
	return err
}
