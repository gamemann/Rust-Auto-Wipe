package processor

import (
	"errors"
	"strings"
)

type WipeData struct {
	WipeDay uint8 	// 0 - 6 (Monday -> Sunday).
	WipeHour uint8	// 0 - 24.
	WipeMin uint8	// 0 - 60.

	TimeZone string

	DeleteMap bool
	DeleteBP bool
	DeletePD bool

	ChangeMapSeeds bool
	MapSeeds []int
	MapSeedPickType uint
	MapSeedsMerge bool 
	NextMapSeed int

	ChangeHostName bool
	HostName string
	NextHostName string

	ChatMsgEnable bool
	ChatMsgAmount uint
	ChatMsg string
}

func (wipedata *WipeData) ProcessData(cfg *config.Config, idx int) error {
	var srv *config.Server
	srv = &cfg.Servers[idx]

	// Make sure we have a valid server. This should never be the case since the array is preallocated to my understanding (and therefore never nil).
	if srv == nil {
		return errors.new("Could not find server at index.")
	}

	// Parse wipe time.
	day := strings.split(srv.)

	return nil
}