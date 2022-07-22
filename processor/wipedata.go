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
	NextMapSeed     int

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

	if len(srv.Timezone) > 0 {
		timezone = srv.Timezone
	}

	wipedata.TimeZone = timezone

	// Check for wipe time override.
	wipetime := cfg.DefaultWipeTime

	if len(srv.WipeTime) > 0 {
		wipetime = srv.WipeTime
	}

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
