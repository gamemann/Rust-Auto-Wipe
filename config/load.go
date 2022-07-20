package config

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
