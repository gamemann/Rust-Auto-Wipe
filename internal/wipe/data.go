package wipe

import (
	"errors"
	"reflect"

	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
)

type Internal struct {
	LatestVersion uint64
}

type WarningMessage struct {
	WarningTime uint
	Message     string
}

type Data struct {
	Enabled    bool
	APIURL     string
	APIToken   string
	DebugLevel int

	PathToServerFiles string

	TimeZone  string
	CronStr   []string
	CronMerge bool

	DeleteMap        bool
	DeleteBP         bool
	DeleteDeaths     bool
	DeleteStates     bool
	DeleteIdentities bool
	DeleteTokens     bool

	DeleteSv bool

	ChangeMapSeeds  bool
	MapSeeds        []int
	MapSeedPickType uint
	MapSeedsMerge   bool

	NextMapSeed int

	ChangeHostName bool
	HostName       string
	NextHostName   string

	MergeWarnings   bool
	WarningMessages []WarningMessage

	InternalData Internal
}

func ProcessData(data *Data, cfg *config.Config, srv *config.Server) error {
	// Make sure we have a valid server. This should never be the case since the array is preallocated to my understanding (and therefore never nil).
	if srv == nil {
		return errors.New("Could not find server at index.")
	}

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

	// Check for time zone override.
	timezone := cfg.Timezone

	if srv.Timezone != nil && len(*srv.Timezone) > 0 {
		timezone = *srv.Timezone
	}

	data.TimeZone = timezone

	// Check for cron merge override.
	cron_merge := cfg.CronMerge

	if srv.CronMerge != nil {
		cron_merge = *srv.CronMerge
	}

	data.CronMerge = cron_merge

	// Check for wipe time override.
	var crons []string

	cron_str := cfg.CronStr

	// Add defaults to cron string slice.
	s := reflect.ValueOf(cron_str)

	if s.Kind() == reflect.String {
		crons = append(crons, s.String())
	} else if s.Kind() == reflect.Slice {
		for i := 0; i < s.Len(); i++ {
			new_cron := s.Index(i).Interface().(string)

			crons = append(crons, new_cron)
		}
	}

	if srv.CronStr != nil {
		s = reflect.ValueOf(*srv.CronStr)

		// Check if string.
		if s.Kind() == reflect.String {
			tmp := s.String()

			if data.CronMerge {
				crons = append(crons, tmp)
			} else {
				crons = []string{tmp}
			}
		} else if s.Kind() == reflect.Slice {
			for i := 0; i < s.Len(); i++ {
				new_cron := s.Index(i).Interface().(string)

				if !data.CronMerge {
					crons = []string{}
				}

				crons = append(crons, new_cron)
			}
		}
	}

	data.CronStr = crons

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

	// Check for delete deaths override.
	delete_deaths := cfg.DeleteDeaths

	if srv.DeleteDeaths != nil {
		delete_deaths = *srv.DeleteDeaths
	}

	data.DeleteDeaths = delete_deaths

	// Check for delete states override.
	delete_states := cfg.DeleteStates

	if srv.DeleteStates != nil {
		delete_states = *srv.DeleteStates
	}

	data.DeleteStates = delete_states

	// Check for delete identities override.
	delete_identities := cfg.DeleteIdentities

	if srv.DeleteIdentities != nil {
		delete_identities = *srv.DeleteIdentities
	}

	data.DeleteIdentities = delete_identities

	// Check for delete tokens override.
	delete_tokens := cfg.DeleteTokens

	if srv.DeleteTokens != nil {
		delete_tokens = *srv.DeleteTokens
	}

	data.DeleteTokens = delete_tokens

	// Check for delete player data override.
	deletesv := cfg.DeleteSv

	if srv.DeleteSv != nil {
		deletesv = *srv.DeleteSv
	}

	data.DeleteSv = deletesv

	// Check for change map seed override.
	changemapseeds := cfg.ChangeMapSeed

	if srv.ChangeMapSeeds != nil {
		changemapseeds = *srv.ChangeMapSeeds
	}

	data.ChangeMapSeeds = changemapseeds

	// Check for map seeds merge in server-specific settings.
	mapseedsmerge := false

	if srv.MapSeedsMerge != nil {
		mapseedsmerge = *srv.MapSeedsMerge
	}

	data.MapSeedsMerge = mapseedsmerge

	// Check for map seeds override.
	var seeds []int

	map_seeds := cfg.MapSeeds

	s = reflect.ValueOf(map_seeds)

	// Check types.
	if s.Kind() == reflect.Int {
		seeds = append(seeds, int(s.Int()))
	} else if s.Kind() == reflect.Slice {
		for i := 0; i < s.Len(); i++ {
			new_seed := int(s.Index(i).Interface().(float64))

			seeds = append(seeds, new_seed)
		}
	}

	if srv.MapSeeds != nil {
		s = reflect.ValueOf(*srv.MapSeeds)

		if s.Kind() == reflect.Int {
			if !data.MapSeedsMerge {
				seeds = []int{}
			} else {
				seeds = []int{int(s.Int())}
			}
		} else if s.Kind() == reflect.Slice {
			if !data.MapSeedsMerge {
				seeds = []int{}
			}

			for i := 0; i < s.Len(); i++ {
				new_seed := int(s.Index(i).Interface().(float64))

				seeds = append(seeds, new_seed)
			}
		}
	}

	data.MapSeeds = seeds

	// Check for map seeds pick type override.
	mapseedspicktype := cfg.MapSeedsPickType

	if srv.MapSeedsPickType != nil {
		mapseedspicktype = *srv.MapSeedsPickType
	}

	data.MapSeedPickType = uint(mapseedspicktype)

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

	// Check for warnings merge override.
	merge_warnings := cfg.MergeWarnings

	if srv.MergeWarnings != nil {
		merge_warnings = *srv.MergeWarnings
	}

	data.MergeWarnings = merge_warnings

	// Check for warnings override.
	var warning_messages []WarningMessage

	// Since this is a custom structure, we have to use somewhat sloppy code :(
	for _, tmp := range cfg.WarningMessages {
		var warning WarningMessage

		warning.WarningTime = tmp.WarningTime
		warning.Message = tmp.Message

		warning_messages = append(warning_messages, warning)
	}

	// Check if we need to merge warning messages or override.
	if srv.WarningMessages != nil {
		if merge_warnings {
			for _, tmp := range *srv.WarningMessages {
				var warning WarningMessage
				warning.Message = tmp.Message
				warning.WarningTime = tmp.WarningTime

				warning_messages = append(warning_messages, warning)
			}
		} else {
			// Wipe messages and override.
			warning_messages = []WarningMessage{}

			for _, tmp := range *srv.WarningMessages {
				var warning WarningMessage
				warning.Message = tmp.Message
				warning.WarningTime = tmp.WarningTime

				warning_messages = append(warning_messages, warning)
			}
		}
	}

	data.WarningMessages = warning_messages

	return nil
}
