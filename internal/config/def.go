package config

type WarningMessage struct {
	WarningTime uint   `json:"warningtime"`
	Message     string `json:"message"`
}

type Server struct {
	// Server ID from Pterodactyl.
	UUID string `json:"uuid"`

	// API/Debug.
	APIURL     *string `json:"apiurl"`
	APIToken   *string `json:"apitoken"`
	DebugLevel *int    `json:"debuglevel"`

	// Paths (e.g. /home/container/server/rust).
	PathToServerFiles *string `json:"pathtoserverfiles"`

	// Wipe date/times.
	Timezone     *string `json:"timezone"`
	WipeTime     *string `json:"wipetime"`
	WipeMonthly  *bool   `json:"wipemonthly"`
	WipeBiweekly *bool   `json:"wipebiweekly"`

	// Files/data that should be deleted.
	DeleteMap        *bool `json:"deletemap"`
	DeleteBP         *bool `json:"deletebp"`
	DeleteDeaths     *bool `json:"deletedeaths"`
	DeleteIdentities *bool `json:"deleteidentities"`
	DeleteStates     *bool `json:"deletestates"`
	DeleteTokens     *bool `json:"deletetokens"`

	DeleteSv *bool `json:"deletesv"`

	// Map seeds.
	ChangeMapSeeds   *bool  `json:"changemapseed"`
	MapSeeds         *[]int `json:"mapseeds"`
	MapSeedsPickType *int   `json:"mapseedspicktype"`
	MapSeedsMerge    *bool  `json:"mergeseeds"`

	// Host name.
	ChangeHostName *bool   `json:"changehostname"`
	HostName       *string `json:"hostname"`

	// Warning chat messages.
	MergeWarnings   *bool             `json:"mergewarnings"`
	WarningMessages *[]WarningMessage `json:"warningmessages"`
}

type Config struct {
	// Pterodactyl API (there are overrides for these).
	APIURL     string `json:"apiurl"`
	APIToken   string `json:"apitoken"`
	DebugLevel int    `json:"debuglevel"`

	// Auto add servers.
	AppToken       string `json:"apptoken"`
	AutoAddServers bool   `json:"autoaddservers"`

	// Paths (e.g. /home/container/server/rust).
	PathToServerFiles string `json:"pathtoserverfiles"`

	// Wipe date times.
	Timezone     string `json:"timezone"`
	WipeTime     string `json:"wipetime"`
	WipeMonthly  bool   `json:"wipemonthly"`
	WipeBiweekly bool   `json:"wipebiweekly"`

	// Files/data that should be deleted.
	DeleteMap        bool `json:"deletemap"`
	DeleteBP         bool `json:"deletebp"`
	DeleteDeaths     bool `json:"deletedeaths"`
	DeleteStates     bool `json:"deletestates"`
	DeleteIdentities bool `json:"deleteidentities"`
	DeleteTokens     bool `json:"deletetokens"`

	DeleteSv bool `json:"deletesv"`

	// Map seeds.
	ChangeMapSeed    bool  `json:"changemapseed"`
	MapSeeds         []int `json:"mapseeds"`
	MapSeedsPickType int   `json:"mapspicktype"`

	// Host name.
	ChangeHostName bool   `json:"changehostname"`
	HostName       string `json:"hostname"`

	// Warning chat messages.
	MergeWarnings   bool             `json:"mergewarnings"`
	WarningMessages []WarningMessage `json:"warningmessages"`

	Servers []Server `json:"servers"`
}

// Sets config's  values.
func (cfg *Config) SetDefaults() {
	cfg.DebugLevel = 1

	cfg.PathToServerFiles = "/server/rust"

	cfg.Timezone = "America/Chicago"
	cfg.WipeTime = "Thursday 12:00"
	cfg.WipeMonthly = false
	cfg.WipeBiweekly = false

	cfg.AutoAddServers = false

	cfg.DeleteMap = true
	cfg.DeleteBP = true
	cfg.DeleteDeaths = true
	cfg.DeleteStates = true
	cfg.DeleteIdentities = true
	cfg.DeleteTokens = true

	cfg.DeleteSv = true

	cfg.ChangeMapSeed = false
	cfg.MapSeedsPickType = 1

	cfg.ChangeHostName = true
	cfg.HostName = "Vanilla | FULL WIPE {month}/{day}"

	// Warn each second for the last 10 seconds before the wipe.
	for i := 1; i <= 10; i++ {
		var warning WarningMessage

		warning.Message = "Wiping server in {seconds_left} seconds. Please join back!"
		warning.WarningTime = uint(i)

		cfg.WarningMessages = append(cfg.WarningMessages, warning)
	}
}
