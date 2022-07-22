package processor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gamemann/Rust-Auto-Wipe/config"
)

type WipeData struct {
	WipeDay  uint8 // 0 - 6 (Monday -> Sunday).
	WipeHour uint8 // 0 - 24.
	WipeMin  uint8 // 0 - 60.

	TimeZone string

	DeleteMap bool
	DeleteBP  bool
	DeletePD  bool

	ChangeMapSeeds  bool
	MapSeeds        []int
	MapSeedPickType uint
	MapSeedsMerge   bool

	NextMapSeed int

	ChangeHostName bool
	HostName       string
	NextHostName   string

	ChatMsgEnable bool
	ChatMsgAmount uint
	ChatMsg       string
}

func (wipedata *WipeData) ProcessData(cfg *config.Config, idx int) error {
	var srv *config.Server
	srv = &cfg.Servers[idx]

	// Make sure we have a valid server. This should never be the case since the array is preallocated to my understanding (and therefore never nil).
	if srv == nil {
		return errors.New("Could not find server at index.")
	}

	// Check for time zone override.
	timezone := cfg.DefaultTimezone

	if srv.Timezone != nil && len(*srv.Timezone) > 0 {
		timezone = *srv.Timezone
	}

	wipedata.TimeZone = timezone

	// Check for wipe time override.
	wipetime := cfg.DefaultWipeTime

	if srv.WipeTime != nil && len(*srv.WipeTime) > 0 {
		wipetime = *srv.WipeTime
	}

	// Check for delete map override.
	deletemap := cfg.DefaultDeleteMap

	if srv.DeleteMap != nil {
		deletemap = *srv.DeleteMap
	}

	wipedata.DeleteMap = deletemap

	// Check for delete blueprint override.
	deletebp := cfg.DefaultDeleteBP

	if srv.DeleteBP != nil {
		deletebp = *srv.DeleteBP
	}

	wipedata.DeleteBP = deletebp

	// Check for delete player data override.
	deletepd := cfg.DefaultDeletePD

	if srv.DeletePD != nil {
		deletepd = *srv.DeletePD
	}

	wipedata.DeletePD = deletepd

	// Check for change map seed override.
	changemapseeds := cfg.DefaultChangeMapSeed

	if srv.ChangeMapSeeds != nil {
		changemapseeds = *srv.ChangeMapSeeds
	}

	wipedata.ChangeMapSeeds = changemapseeds

	// Check for map seeds override.
	mapseeds := cfg.DefaultMapSeeds

	if srv.MapSeeds != nil {
		mapseeds = *srv.MapSeeds
	}

	wipedata.MapSeeds = mapseeds

	// Check for map seeds pick type override.
	mapseedspicktype := cfg.DefaultMapSeedsPickType

	if srv.MapSeedsPickType != nil {
		mapseedspicktype = *srv.MapSeedsPickType
	}

	wipedata.MapSeedPickType = uint(mapseedspicktype)

	// Check for map seeds merge in server-specific settings.
	mapseedsmerge := false

	if srv.MapSeedsMerge != nil {
		mapseedsmerge = *srv.MapSeedsMerge
	}

	wipedata.MapSeedsMerge = mapseedsmerge

	// Check for change host name override.
	changehostname := cfg.DefaultChangeHostName

	if srv.ChangeHostName != nil {
		changehostname = *srv.ChangeHostName
	}

	wipedata.ChangeHostName = changehostname

	// Check for host name override.
	hostname := cfg.DefaultHostName

	if srv.HostName != nil {
		hostname = *srv.HostName
	}

	wipedata.HostName = hostname

	// Check for chat message enable override.
	chatmsgenable := cfg.DefaultChatMsgEnable

	if srv.ChatMsgEnable != nil {
		chatmsgenable = *srv.ChatMsgEnable
	}

	wipedata.ChatMsgEnable = chatmsgenable

	// Check for chat message amount override.
	chatmsgamount := cfg.DefaultChatMsgAmount

	if srv.ChatMsgAmount != nil {
		chatmsgamount = *srv.ChatMsgAmount
	}

	wipedata.ChatMsgAmount = uint(chatmsgamount)

	// Check for chat message override.
	chatmsg := cfg.DefaultChatMsg

	if srv.ChatMsg != nil {
		chatmsg = *srv.ChatMsg
	}

	wipedata.ChatMsg = chatmsg

	// Parse wipe time.
	info := strings.Split(wipetime, " ")

	day := info[0]
	timeinfo := info[1]

	// Convert day to numberic value from 0 - 6.
	switch strings.ToLower(day) {
	case "monday":
		wipedata.WipeDay = 0
	case "tuesday":
		wipedata.WipeDay = 1
	case "wednesday":
		wipedata.WipeDay = 2
	case "thursday":
		wipedata.WipeDay = 3
	case "friday":
		wipedata.WipeDay = 4
	case "saturday":
		wipedata.WipeDay = 5
	case "sunday":
		wipedata.WipeDay = 6
	}

	td := strings.Split(timeinfo, ":")

	hour, err := strconv.Atoi(td[0])

	if err != nil {
		return err
	}

	min, err := strconv.Atoi(td[1])

	// Do boundary checks.
	if hour > 24 {
		hour = 24

		if cfg.DebugLevel > 0 {
			fmt.Println("[WARNING] Found hour out of bounds. (> 24).")
		}
	} else if hour < 0 {
		hour = 0

		if cfg.DebugLevel > 0 {
			fmt.Println("[WARNING] Found hour out of bounds (< 0).")
		}
	}

	if min > 60 {
		min = 60

		if cfg.DebugLevel > 0 {
			fmt.Println("[WARNING] Found minute out of bounds (> 60).")
		}
	} else if min < 0 {
		min = 0

		if cfg.DebugLevel > 0 {
			fmt.Println("[WARNING] Found minute out of bounds (< 0).")
		}
	}

	return nil
}
