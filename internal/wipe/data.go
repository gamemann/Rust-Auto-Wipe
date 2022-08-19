package wipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

const (
	Default_Map       = "Procedural Map"
	Default_WorldSize = 3000
	Default_WorldSeed = 0

	Default_WarningMessage = "Wiping server in {seconds_left} seconds. Please join back!"
)

type Internal struct {
	LatestVersion uint64
	LatestWorld   uint
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
	DeleteSv         bool
	DeleteFiles      []config.AdditionalFiles
	DeleteFilesMerge bool

	ChangeWorldInfo   bool
	WorldInfo         []config.WorldInfo
	WorldInfoPickType uint
	WorldInfoMerge    bool

	NextMapSeed int

	ChangeHostName bool
	HostName       string
	NextHostName   string

	MergeWarnings   bool
	WarningMessages []config.WarningMessage

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
		if len(*srv.APIURL) > 0 {
			apiurl = *srv.APIURL
		}
	}

	data.APIURL = apiurl

	// Check for API token override.
	apitoken := cfg.APIToken

	if srv.APIToken != nil {
		if len(*srv.APIToken) > 0 {
			apitoken = *srv.APIToken
		}
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
		if len(*srv.PathToServerFiles) > 0 {
			pathtoserverfiles = *srv.PathToServerFiles
		}
	}

	data.PathToServerFiles = pathtoserverfiles

	// Check for time zone override.
	timezone := cfg.Timezone

	if srv.Timezone != nil && len(*srv.Timezone) > 0 {
		if len(*srv.Timezone) > 0 {
			timezone = *srv.Timezone
		}
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

	if s.Kind() == reflect.String && s.Len() > 0 {
		crons = append(crons, s.String())
	} else if s.Kind() == reflect.Slice {
		for i := 0; i < s.Len(); i++ {
			new_cron := s.Index(i).Interface().(string)

			if len(new_cron) < 1 {
				continue
			}

			crons = append(crons, new_cron)
		}
	}

	if srv.CronStr != nil {
		s = reflect.ValueOf(*srv.CronStr)

		// Check if string.
		if s.Kind() == reflect.String && s.Len() > 0 {
			tmp := s.String()

			if data.CronMerge {
				crons = append(crons, tmp)
			} else {
				crons = []string{tmp}
			}
		} else if s.Kind() == reflect.Slice {
			for i := 0; i < s.Len(); i++ {
				new_cron := s.Index(i).Interface().(string)

				if len(new_cron) < 1 {
					continue
				}

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

	// Check for warnings merge override.
	delete_files_merge := cfg.DeleteFilesMerge

	if srv.DeleteFilesMerge != nil {
		delete_files_merge = *srv.DeleteFilesMerge
	}

	data.DeleteFilesMerge = delete_files_merge

	// Check for warnings override.
	delete_files := cfg.DeleteFiles

	// Check if we need to merge warning messages or override.
	if srv.DeleteFiles != nil {
		if delete_files_merge {
			for _, v := range *srv.DeleteFiles {
				delete_files = append(delete_files, v)
			}
		} else {
			delete_files = *srv.DeleteFiles
		}
	}

	data.DeleteFiles = delete_files

	// Check for change world info override.
	changeworldinfo := cfg.ChangeWorldInfo

	if srv.ChangeWorldInfo != nil {
		changeworldinfo = *srv.ChangeWorldInfo
	}

	data.ChangeWorldInfo = changeworldinfo

	// Check for world info merge in server-specific settings.
	worldinfomerge := false

	if srv.WorldInfoMerge != nil {
		worldinfomerge = *srv.WorldInfoMerge
	}

	data.WorldInfoMerge = worldinfomerge

	// Check for world info override.
	world_info := cfg.WorldInfo

	if srv.WorldInfo != nil {
		if data.WorldInfoMerge {
			for _, v := range *srv.WorldInfo {
				world_info = append(world_info, v)
			}
		} else {
			world_info = *srv.WorldInfo
		}
	}

	data.WorldInfo = world_info

	// Check for world seed pick type override.
	world_info_pick_type := cfg.WorldInfoPickType

	if srv.WorldInfoPickType != nil {
		world_info_pick_type = *srv.WorldInfoPickType
	}

	data.WorldInfoPickType = uint(world_info_pick_type)

	// Check for change host name override.
	changehostname := cfg.ChangeHostName

	if srv.ChangeHostName != nil {
		changehostname = *srv.ChangeHostName
	}

	data.ChangeHostName = changehostname

	// Check for host name override.
	hostname := cfg.HostName

	if srv.HostName != nil {
		if len(*srv.HostName) > 0 {
			hostname = *srv.HostName
		}
	}

	data.HostName = hostname

	// Check for warnings merge override.
	merge_warnings := cfg.MergeWarnings

	if srv.MergeWarnings != nil {
		merge_warnings = *srv.MergeWarnings
	}

	data.MergeWarnings = merge_warnings

	// Check for warnings override.
	warning_messages := cfg.WarningMessages

	// Check if we need to merge warning messages or override.
	if srv.WarningMessages != nil {
		if merge_warnings {
			for _, v := range *srv.WarningMessages {
				warning_messages = append(warning_messages, v)
			}
		} else {
			warning_messages = *srv.WarningMessages
		}
	}

	data.WarningMessages = warning_messages

	// Fill out null data for world information.
	ep := "client/servers/" + srv.UUID + "/startup"

	// We first need to retrieve the current variable.
	d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "GET", ep, nil)

	debug.SendDebugMsg(srv.UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => nil.")
	debug.SendDebugMsg(srv.UUID, data.DebugLevel, 4, "List Variable return data => "+d+".")

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(srv.UUID, data.DebugLevel, 0, "Could not list startup variables. Please enable debugging level 4 for body response including errors.")

		return errors.New("Could not list startup variables.")
	}

	if err != nil {
		return err
	}

	// We want to parse the response with the startup response structure.
	var EnvVars pterodactyl.StartupResp

	// Convert to JSON.
	err = json.Unmarshal([]byte(d), &EnvVars)

	if err != nil {
		return err
	}

	for _, seed := range EnvVars.Data {
		if seed.Attributes.Env_Variable == "WORLD_SEED" && len(seed.Attributes.Env_Variable) > 0 {
			s := seed.Attributes.Srv_Value

			for i := 0; i < len(data.WorldInfo); i++ {
				if data.WorldInfo[i].WorldSeed == nil {
					s_int, err := strconv.Atoi(s)

					if err != nil {
						fmt.Println(err)

						continue
					}

					data.WorldInfo[i].WorldSeed = &s_int
				}
			}
		}

		if seed.Attributes.Env_Variable == "WORLD_SIZE" && len(seed.Attributes.Env_Variable) > 0 {
			s := seed.Attributes.Srv_Value

			for i := 0; i < len(data.WorldInfo); i++ {
				if data.WorldInfo[i].WorldSize == nil {
					s_int, err := strconv.Atoi(s)

					if err != nil {
						fmt.Println(err)

						continue
					}

					data.WorldInfo[i].WorldSize = &s_int
				}
			}
		}

		if seed.Attributes.Env_Variable == "LEVEL" && len(seed.Attributes.Env_Variable) > 0 {
			s := seed.Attributes.Srv_Value

			for i := 0; i < len(data.WorldInfo); i++ {
				if data.WorldInfo[i].Map == nil {
					data.WorldInfo[i].Map = &s
				}
			}
		}
	}

	// Loop through world information one more time and fill out anything missed with constants.
	for i := 0; i < len(data.WorldInfo); i++ {
		if data.WorldInfo[i].Map == nil {
			tmp := Default_Map

			data.WorldInfo[i].Map = &tmp
		}

		if data.WorldInfo[i].WorldSize == nil {
			tmp := Default_WorldSeed

			data.WorldInfo[i].WorldSize = &tmp
		}

		if data.WorldInfo[i].WorldSeed == nil {
			tmp := Default_WorldSize

			data.WorldInfo[i].WorldSize = &tmp
		}
	}

	// Loop through warning messages and fill out warning message string if nil.
	for i := 0; i < len(data.WarningMessages); i++ {
		if data.WarningMessages[i].Message == nil {
			tmp := Default_WarningMessage

			data.WarningMessages[i].Message = &tmp
		}
	}

	return nil
}
