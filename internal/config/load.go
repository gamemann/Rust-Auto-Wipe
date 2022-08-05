package config

import (
	"encoding/json"
	"os"
)

// Reads a config file based off of the file name (string) and returns a Config struct.
func (cfg *Config) LoadConfig(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	stat, _ := file.Stat()

	data := make([]byte, stat.Size())

	_, err = file.Read(data)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), cfg)

	return err
}
