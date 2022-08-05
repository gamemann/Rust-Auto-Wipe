package wipe

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
)

type Internal struct {
	LatestVersion uint64
	LastDayNum    int
	LastMonthNum  int
}

type Data struct {
	WipeDay      uint8 // 0 - 6 (Sunday -> Saturday).
	WipeHour     uint8 // 0 - 24.
	WipeMin      uint8 // 0 - 60.
	WipeMonthly  bool
	WipeBiweekly bool

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

	APIURL     string
	APIToken   string
	DebugLevel int

	PathToServerFiles string

	InternalData Internal
}

func ProcessData(data *Data, cfg *config.Config, srv *config.Server) error {
	// Make sure we have a valid server. This should never be the case since the array is preallocated to my understanding (and therefore never nil).
	if srv == nil {
		return errors.New("Could not find server at index.")
	}

	// Check for time zone override.
	timezone := cfg.Timezone

	if srv.Timezone != nil && len(*srv.Timezone) > 0 {
		timezone = *srv.Timezone
	}

	data.TimeZone = timezone

	// Check for wipe time override.
	wipetime := cfg.WipeTime

	if srv.WipeTime != nil && len(*srv.WipeTime) > 0 {
		wipetime = *srv.WipeTime
	}

	// Check for wipe monthly override.
	wipemonthly := cfg.WipeMonthly

	if srv.WipeMonthly != nil {
		wipemonthly = *srv.WipeMonthly
	}

	data.WipeMonthly = wipemonthly

	// Check for wipe biweekly override.
	wipebiweekly := cfg.WipeBiweekly

	if srv.WipeBiweekly != nil {
		wipebiweekly = *srv.WipeBiweekly
	}

	data.WipeBiweekly = wipebiweekly

	// Check for delete map override.
	deletemap := cfg.DeleteMap

	if srv.DeleteMap != nil {
		deletemap = *srv.DeleteMap
	}

	data.DeleteMap = deletemap

	// Check for delete blueprint override.
	deletebp := cfg.DeleteBP

	if srv.DeleteBP != nil {
		deletebp = *srv.DeleteBP
	}

	data.DeleteBP = deletebp

	// Check for delete player data override.
	deletepd := cfg.DeletePD

	if srv.DeletePD != nil {
		deletepd = *srv.DeletePD
	}

	data.DeletePD = deletepd

	// Check for change map seed override.
	changemapseeds := cfg.ChangeMapSeed

	if srv.ChangeMapSeeds != nil {
		changemapseeds = *srv.ChangeMapSeeds
	}

	data.ChangeMapSeeds = changemapseeds

	// Check for map seeds override.
	mapseeds := cfg.MapSeeds

	if srv.MapSeeds != nil {
		mapseeds = *srv.MapSeeds
	}

	data.MapSeeds = mapseeds

	// Check for map seeds pick type override.
	mapseedspicktype := cfg.MapSeedsPickType

	if srv.MapSeedsPickType != nil {
		mapseedspicktype = *srv.MapSeedsPickType
	}

	data.MapSeedPickType = uint(mapseedspicktype)

	// Check for map seeds merge in server-specific settings.
	mapseedsmerge := false

	if srv.MapSeedsMerge != nil {
		mapseedsmerge = *srv.MapSeedsMerge
	}

	data.MapSeedsMerge = mapseedsmerge

	// Check for change host name override.
	changehostname := cfg.ChangeHostName

	if srv.ChangeHostName != nil {
		changehostname = *srv.ChangeHostName
	}

	data.ChangeHostName = changehostname

	// Check for host name override.
	hostname := cfg.HostName

	if srv.HostName != nil {
		hostname = *srv.HostName
	}

	data.HostName = hostname

	// Check for chat message enable override.
	chatmsgenable := cfg.ChatMsgEnable

	if srv.ChatMsgEnable != nil {
		chatmsgenable = *srv.ChatMsgEnable
	}

	data.ChatMsgEnable = chatmsgenable

	// Check for chat message amount override.
	chatmsgamount := cfg.ChatMsgAmount

	if srv.ChatMsgAmount != nil {
		chatmsgamount = *srv.ChatMsgAmount
	}

	data.ChatMsgAmount = uint(chatmsgamount)

	// Check for chat message override.
	chatmsg := cfg.ChatMsg

	if srv.ChatMsg != nil {
		chatmsg = *srv.ChatMsg
	}

	data.ChatMsg = chatmsg

	// Check for API URL override.
	apiurl := cfg.APIURL

	if srv.APIURL != nil {
		apiurl = *srv.APIURL
	}

	data.APIURL = apiurl

	// Check for API token override.
	apitoken := cfg.APIToken

	if srv.APIToken != nil {
		apitoken = *srv.APIToken
	}

	data.APIToken = apitoken

	// Check for debug level override.
	debuglevel := cfg.DebugLevel

	if srv.DebugLevel != nil {
		debuglevel = *srv.DebugLevel
	}

	data.DebugLevel = debuglevel

	// Check for path to server files override.
	pathtoserverfiles := cfg.PathToServerFiles

	if srv.PathToServerFiles != nil {
		pathtoserverfiles = *srv.PathToServerFiles
	}

	data.PathToServerFiles = pathtoserverfiles

	// Parse wipe time.
	info := strings.Split(wipetime, " ")

	day := info[0]
	timeinfo := info[1]

	// Convert day to numberic value from 0 - 6.
	switch strings.ToLower(day) {
	case "sunday":
		data.WipeDay = 0
	case "monday":
		data.WipeDay = 1
	case "tuesday":
		data.WipeDay = 2
	case "wednesday":
		data.WipeDay = 3
	case "thursday":
		data.WipeDay = 4
	case "friday":
		data.WipeDay = 5
	case "saturday":
		data.WipeDay = 6
	default:
		data.WipeDay = 0
	}

	td := strings.Split(timeinfo, ":")

	// Make sure we have two or more elements.
	if len(td) < 2 {
		return errors.New("Time info split failure (< 2 array size).")
	}

	hour, err := strconv.Atoi(td[0])

	if err != nil {
		return err
	}

	min, err := strconv.Atoi(td[1])

	if err != nil {
		return err
	}

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
