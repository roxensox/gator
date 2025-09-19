package config

import (
	"encoding/json"
	"os"
)

func (c *Config) Read() error {
	// Reads the config JSON file into a Config object

	// Gets the path to the config file
	configPath, err := getConfigFilePath()
	// If there is an error getting the path, returns it
	if err != nil {
		return err
	}

	// Reads the file into a byte array
	data, err := os.ReadFile(configPath)
	// If there is an error reading the file, returns it
	if err != nil {
		return err
	}

	// Unloads the file data into the Config object
	json.Unmarshal(data, c)
	return nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUser = &user
	return write(*c)
}

func NewConfig() Config {
	// Creates and returns a new Config object

	c := Config{}
	return c
}
