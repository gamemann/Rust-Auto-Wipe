package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Server struct {
	UUID int `json:"uuid"`

	// Wipe date/times.
	Timezone string `json:"timezone"`
	WipeTime string `json:"wipetime"`

	// Files/data that should be deleted.
	DeleteMap bool `json:"deletemap"`
	DeleteBP  bool `json:"deletebp"`
	DeletePD  bool `json:"deletepd"`

	// Map seeds.
	ChangeMapSeeds   bool  `json:"changemapseed"`
	MapSeeds         []int `json:"mapseeds"`
	MapSeedsPickType int   `json:"mapseedspicktype"`

	// Host name.
	ChangeHostName bool   `json:"changehostname"`
	HostName       string `json:"hostname"`

	// Warning chat messages.
	ChatMsgEnable bool   `json:"chatmsgenable"`
	ChatMsg       string `json:"chatmsg"`
	ChatMsgAmount int    `json:"chatmsgamount"`
}

type Config struct {
	// Pterodactyl API.
	APIEndpoint string `json:"apiurl"`
	APIToken    string `json:"apitoken"`

	// Paths (e.g. /home/container/server/rust).
	DefaultPathToServerFiles string `json:"defaultpathtoserverfiles"`

	// Wipe date times.
	DefaultTimezone string `json:"defaulttimezone"`
	DefaultWipeTime string `json:"defaultwipetime"`

	// Files/data that should be deleted.
	DefaultDeleteMap bool `json:"defaultdeletemap"`
	DefaultDeleteBP  bool `json:"defaultdeletebp"`
	DefaultDeletePD  bool `json:"defaultdeletepd"`

	// Map seeds.
	DefaultChangeMapSeed    bool  `json:"defaultchangemapseed"`
	DefaultMapSeeds         []int `json:"defaultmapseeds"`
	DefaultMapSeedsPickType int   `json:"defaultmapspicktype"`

	// Host name.
	DefaultChangeHostName bool   `json:"defaultchangehostname"`
	DefaultHostName       string `json:"defaulthostname"`

	// Warning chat messages.
	DefaultChatMsgEnable bool   `json:"defaultchatmsgenable"`
	DefaultChatMsg       string `json:"defaultchatmsg"`
	DefaultChatMsgAmount int    `json:"defaultchatmsgamount"`

	Servers []Server `json:"servers"`
}

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

// Sets config's default values.
func (cfg *Config) SetDefaults() {
	cfg.DefaultPathToServerFiles = "/home/container/server/rust"

	cfg.DefaultTimezone = "America/Chicago"
	cfg.DefaultWipeTime = ""

	cfg.DefaultDeleteMap = true
	cfg.DefaultDeleteBP = true
	cfg.DefaultDeletePD = true

	cfg.DefaultChangeMapSeed = false
	cfg.DefaultMapSeedsPickType = 1

	cfg.DefaultChangeHostName = true
	cfg.DefaultHostName = "Vanilla | FULLWIPE {wipetime_one}"

	cfg.DefaultChatMsgEnable = true
	cfg.DefaultChatMsgAmount = 5
	cfg.DefaultChatMsg = "Wiping server in {seconds} seconds. Please join back!"
}
