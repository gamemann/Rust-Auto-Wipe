package config

type WarningMessage struct {
	WarningTime uint    `json:"warningtime"`
	Message     *string `json:"message"`
}

type WorldInfo struct {
	Map       *string `json:"map"`
	WorldSize *int    `json:"worldsize"`
	WorldSeed *int    `json:"worldseed"`
}

type AdditionalFiles struct {
	Root  string   `json:"root"`
	Files []string `json:"files"`
}

type Server struct {
	Enabled bool `json:"enabled"`

	// Server ID from Pterodactyl.
	ID     int    `json:"id"`
	LongID string `json:"uuidlong"`
	UUID   string `json:"uuid"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Name   string `json:"name"`

	// API/Debug.
	APIURL     *string `json:"apiurl"`
	APIToken   *string `json:"apitoken"`
	DebugLevel *int    `json:"debuglevel"`

	// Paths (e.g. /server/rust).
	PathToServerFiles *string `json:"pathtoserverfiles"`

	// Wipe date/times.
	Timezone  *string      `json:"timezone"`
	CronStr   *interface{} `json:"cronstr"`
	CronMerge *bool        `json:"cronmerge"`

	// Files/data that should be deleted.
	DeleteMap        *bool              `json:"deletemap"`
	DeleteBP         *bool              `json:"deletebp"`
	DeleteDeaths     *bool              `json:"deletedeaths"`
	DeleteIdentities *bool              `json:"deleteidentities"`
	DeleteStates     *bool              `json:"deletestates"`
	DeleteTokens     *bool              `json:"deletetokens"`
	DeleteSv         *bool              `json:"deletesv"`
	DeleteFiles      *[]AdditionalFiles `json:"deletefiles"`
	DeleteFilesMerge *bool              `json:"deletefilesmerge"`

	// Map seeds/sizes.
	ChangeWorldInfo   *bool        `json:"changeworldinfo"`
	WorldInfo         *[]WorldInfo `json:"worldinfo"`
	WorldInfoPickType *int         `json:"worldinfopicktype"`
	WorldInfoMerge    *bool        `json:"worldinfomerge"`

	// Host name.
	ChangeHostName *bool   `json:"changehostname"`
	HostName       *string `json:"hostname"`

	// Warning chat messages.
	MergeWarnings   *bool             `json:"mergewarnings"`
	WarningMessages *[]WarningMessage `json:"warningmessages"`

	// Extras (e.g. development testing, etc.).
	WipeFirst bool `json:"wipefirst"`
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
	Timezone  string      `json:"timezone"`
	CronStr   interface{} `json:"cronstr"`
	CronMerge bool        `json:"cronmerge"`

	// Files/data that should be deleted.
	DeleteMap        bool              `json:"deletemap"`
	DeleteBP         bool              `json:"deletebp"`
	DeleteDeaths     bool              `json:"deletedeaths"`
	DeleteStates     bool              `json:"deletestates"`
	DeleteIdentities bool              `json:"deleteidentities"`
	DeleteTokens     bool              `json:"deletetokens"`
	DeleteSv         bool              `json:"deletesv"`
	DeleteFiles      []AdditionalFiles `json:"deletefiles"`
	DeleteFilesMerge bool              `json:"deletefilesmerge"`

	// Maps, seeds, and world sizes.
	ChangeWorldInfo   bool        `json:"changeworldinfo"`
	WorldInfo         []WorldInfo `json:"worldinfo"`
	WorldInfoPickType int         `json:"worldinfopicktype"`
	WorldInfoMerge    bool        `json:"worldinfomerge"`

	// Host name.
	ChangeHostName bool   `json:"changehostname"`
	HostName       string `json:"hostname"`

	// Warning chat messages.
	MergeWarnings   bool             `json:"mergewarnings"`
	WarningMessages []WarningMessage `json:"warningmessages"`

	// Hooks that send information to an endpoint with POST data. Useful for sending updates to Discord for example.
	// Pre wipe hook.
	PreHookAuth string `json:"prehookauth"`
	PreHook     string `json:"prehook"`

	// POST wipe hook.
	PostHookAuth string `json:"posthookauth"`
	PostHook     string `json:"posthook"`

	Servers []Server `json:"servers"`
}

// Sets config's  values.
func (cfg *Config) SetDefaults() {
	cfg.DebugLevel = 1

	cfg.PathToServerFiles = "/server/rust"

	cfg.Timezone = "America/Chicago"
	cfg.CronStr = "30 15 * * 4"
	cfg.CronMerge = true

	cfg.DeleteMap = true
	cfg.DeleteBP = true
	cfg.DeleteDeaths = true
	cfg.DeleteStates = true
	cfg.DeleteIdentities = true
	cfg.DeleteTokens = true

	cfg.DeleteSv = true

	cfg.WorldInfoPickType = 1

	cfg.ChangeHostName = true
	cfg.HostName = "Vanilla | FULL WIPE {month_two}/{day_two}"

	// Warn each second for the last 10 seconds before the wipe.
	for i := 1; i <= 10; i++ {
		var warning WarningMessage
		tmp := "Wiping server in {seconds_left} seconds. Please join back!"

		warning.Message = &tmp

		*warning.Message = "Wiping server in {seconds_left} seconds. Please join back!"
		warning.WarningTime = uint(i)

		cfg.WarningMessages = append(cfg.WarningMessages, warning)
	}
}
