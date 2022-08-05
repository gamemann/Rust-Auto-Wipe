package config

import (
	"encoding/json"
	"os"
)

func (cfg *Config) WriteDefaultsToFile(file string) error {
	var err error

	err = os.MkdirAll("/etc/raw", 0755)

	// If we have an error and it doesn't look like an "already exist" error, return the error.
	if err != nil && !os.IsExist(err) {
		return err
	}

	fp, err := os.Create(file)

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "   ")

	if err != nil {
		// Close file.
		fp.Close()

		return err
	}

	_, err = fp.Write(data)

	if err != nil {
		// Close file.
		fp.Close()

		return err
	}

	return err
}
