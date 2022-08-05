package config

import (
	"encoding/json"
	"os"
)

func (cfg *Config) WriteDefaultsToFile(file string) error {
	var err error

	fp, err := os.Create(file)

	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)

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
