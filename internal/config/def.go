package config

type Server struct {
	// Server ID from Pterodactyl.
	UUID string `json:"uuid"`

	// Wipe date/times.
	Timezone     *string `json:"timezone"`
	WipeTime     *string `json:"wipetime"`
	WipeMonthly  *bool   `json:"wipemonthly"`
	WipeBiweekly *bool   `json:"wipebiweekly"`

	// Files/data that should be deleted.
	DeleteMap *bool `json:"deletemap"`
	DeleteBP  *bool `json:"deletebp"`
	DeletePD  *bool `json:"deletepd"`

	// Map seeds.
	ChangeMapSeeds   *bool  `json:"changemapseed"`
	MapSeeds         *[]int `json:"mapseeds"`
	MapSeedsPickType *int   `json:"mapseedspicktype"`
	MapSeedsMerge    *bool  `json:"mergeseeds"`

	// Host name.
	ChangeHostName *bool   `json:"changehostname"`
	HostName       *string `json:"hostname"`

	// Warning chat messages.
	ChatMsgEnable *bool   `json:"chatmsgenable"`
	ChatMsg       *string `json:"chatmsg"`
	ChatMsgAmount *int    `json:"chatmsgamount"`

	// API/Debug.
	APIURL     *string `json:"apiurl"`
	APIToken   *string `json:"apitoken"`
	DebugLevel *int    `json:"debuglevel"`
}

type Config struct {
	// Pterodactyl API (there are overrides for these).
	APIURL     string `json:"apiurl"`
	APIToken   string `json:"apitoken"`
	DebugLevel int    `json:"debuglevel"`

	// Paths (e.g. /home/container/server/rust).
	PathToServerFiles string `json:"pathtoserverfiles"`

	// Wipe date times.
	Timezone     string `json:"timezone"`
	WipeTime     string `json:"wipetime"`
	WipeMonthly  bool   `json:"wipemonthly"`
	WipeBiweekly bool   `json:"wipebiweekly"`

	// Files/data that should be deleted.
	DeleteMap bool `json:"deletemap"`
	DeleteBP  bool `json:"deletebp"`
	DeletePD  bool `json:"deletepd"`

	// Map seeds.
	ChangeMapSeed    bool  `json:"changemapseed"`
	MapSeeds         []int `json:"mapseeds"`
	MapSeedsPickType int   `json:"mapspicktype"`

	// Host name.
	ChangeHostName bool   `json:"changehostname"`
	HostName       string `json:"hostname"`

	// Warning chat messages.
	ChatMsgEnable bool   `json:"chatmsgenable"`
	ChatMsg       string `json:"chatmsg"`
	ChatMsgAmount int    `json:"chatmsgamount"`

	Servers []Server `json:"servers"`
}

// Sets config's  values.
func (cfg *Config) SetDefaults() {
	cfg.DebugLevel = 1

	cfg.PathToServerFiles = "/home/container/server/rust"

	cfg.Timezone = "America/Chicago"
	cfg.WipeTime = "Thursday 12:00"
	cfg.WipeMonthly = false
	cfg.WipeBiweekly = false

	cfg.DeleteMap = true
	cfg.DeleteBP = true
	cfg.DeletePD = true

	cfg.ChangeMapSeed = false
	cfg.MapSeedsPickType = 1

	cfg.ChangeHostName = true
	cfg.HostName = "Vanilla | FULL WIPE {wipetime_one}"

	cfg.ChatMsgEnable = true
	cfg.ChatMsgAmount = 5
	cfg.ChatMsg = "Wiping server in {seconds} seconds. Please join back!"
}
