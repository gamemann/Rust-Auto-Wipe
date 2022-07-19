package config

type Server struct {
	UUID          int    `json:"uuid"`
	Timezone      string `json:"timezone"`
	WipeTime      string `json:"wipetime"`
	DeleteMap     bool   `json:"deletemap"`
	DeleteBP      bool   `json:"deletebp"`
	DeletePD      bool   `json:"deletepd"`
	MapSeeds      []int  `json:"mapseeds"`
	ChatMsg       string `json:"chatmsg"`
	ChatMsgAmount int    `json:"chatmsgamount"`
}

type Config struct {
	DefaultTimezone      string `json:"defaulttimezone"`
	DefaultWipeTime      string `json:"defaultwipetime"`
	DefaultDeleteMap     bool   `json:"defaultdeletemap"`
	DefaultDeleteBP      bool   `json:"defaultdeletebp"`
	DefaultDeletePD      bool   `json:"defaultdeletepd"`
	DefaultMapSeeds      []int  `json:"defaultmapseeds"`
	DefaultChatMsg       string `json:"defaultchatmsg"`
	DefaultChatMsgAmount int    `json:"defaultchatmsgamount"`

	Servers []Server `json:"servers"`
}
