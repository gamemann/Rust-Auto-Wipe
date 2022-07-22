package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Reads a config file based off of the file name (string) and returns a Config struct.
func (cfg *Config) LoadConfig(path string) bool {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("[ERR] Cannot open config file.")
		fmt.Println(err)

		return false
	}

	defer file.Close()

	stat, _ := file.Stat()

	data := make([]byte, stat.Size())

	_, err = file.Read(data)

	if err != nil {
		fmt.Println("[ERR] Cannot read config file.")
		fmt.Println(err)

		return false
	}

	err = json.Unmarshal([]byte(data), cfg)

	if err != nil {
		fmt.Println("[ERR] Cannot parse JSON Data.")
		fmt.Println(err)

		return false
	}

	return true
}