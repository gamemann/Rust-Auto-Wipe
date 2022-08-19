package wipe

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/internal/autoaddservers"
	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

type ServerDetails struct {
	Object     string `json:"object"`
	Attributes struct {
		ID          int         `json:"id"`
		ExternalID  interface{} `json:"external_id"`
		UUID        string      `json:"uuid"`
		Identifier  string      `json:"identifier"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Status      interface{} `json:"status"`
		Suspended   bool        `json:"suspended"`
		Limits      struct {
			Memory      int         `json:"memory"`
			Swap        int         `json:"swap"`
			Disk        int         `json:"disk"`
			Io          int         `json:"io"`
			CPU         int         `json:"cpu"`
			Threads     interface{} `json:"threads"`
			OomDisabled bool        `json:"oom_disabled"`
		} `json:"limits"`
		FeatureLimits struct {
			Databases   int `json:"databases"`
			Allocations int `json:"allocations"`
			Backups     int `json:"backups"`
		} `json:"feature_limits"`
		User       int `json:"user"`
		Node       int `json:"node"`
		Allocation int `json:"allocation"`
		Nest       int `json:"nest"`
		Egg        int `json:"egg"`
		Container  struct {
			StartupCommand string                `json:"startup_command"`
			Image          string                `json:"image"`
			Installed      int                   `json:"installed"`
			Environment    autoaddservers.RawEnv `json:"environment"`
		} `json:"container"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"attributes"`
}

func EnvOverride(cfg *config.Config, srv *config.Server) error {
	var err error

	ep := "application/servers/" + strconv.Itoa(srv.ID)

	d, _, err := pterodactyl.SendAPIRequest(cfg.APIURL, cfg.AppToken, "GET", ep, nil)

	debug.SendDebugMsg(srv.UUID, cfg.DebugLevel, 3, "Sending request. Request => "+cfg.APIURL+"api/"+ep+". Post data => nil.")
	debug.SendDebugMsg(srv.UUID, cfg.DebugLevel, 4, "List Server Details return data => "+d+".")

	if err != nil {
		return err
	}

	// Convert JSON to structure.
	var server_details ServerDetails

	err = json.Unmarshal([]byte(d), &server_details)

	if err != nil {
		return err
	}

	env := server_details.Attributes.Container.Environment

	// Enabled override.
	if env.RAW_Enabled != nil && len(*env.RAW_Enabled) > 0 {
		srv.Enabled, _ = strconv.ParseBool(*env.RAW_Enabled)
	}

	// Path to server files override.
	if env.RAW_PathToServerFiles != nil && len(*env.RAW_PathToServerFiles) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.PathToServerFiles == nil {
			var val string
			srv.PathToServerFiles = &val
		}

		*srv.PathToServerFiles = *env.RAW_PathToServerFiles
	}

	// Timezone override.
	if env.RAW_Timezone != nil && len(*env.RAW_Timezone) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.Timezone == nil {
			var val string
			srv.Timezone = &val
		}

		*srv.Timezone = *env.RAW_Timezone
	}

	// Cron merge override.
	if env.RAW_CronMerge != nil && len(*env.RAW_CronMerge) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.CronMerge == nil {
			var val bool
			srv.CronMerge = &val
		}

		*srv.CronMerge, _ = strconv.ParseBool(*env.RAW_CronMerge)
	}

	// Cron string override.
	if env.RAW_CronStr != nil && len(*env.RAW_CronStr) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.CronStr == nil {
			var val interface{}
			srv.CronStr = &val
		}

		s := *env.RAW_CronStr
		var tmp []string

		// Try to parse as JSON, if fails, parse as string.
		err := json.Unmarshal([]byte(s), &tmp)

		var str interface{} = ""

		if srv.CronStr == nil {
			srv.CronStr = &str
		}

		if err != nil {
			*srv.CronStr = *env.RAW_CronStr
		} else {
			*srv.CronStr = tmp
		}
	}

	// Delete map override.
	if env.RAW_DeleteMap != nil && len(*env.RAW_DeleteMap) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.DeleteMap == nil {
			var val bool
			srv.DeleteMap = &val
		}

		*srv.DeleteMap, _ = strconv.ParseBool(*env.RAW_DeleteMap)
	}

	// Delete blueprints override.
	if env.RAW_DeleteBP != nil && len(*env.RAW_DeleteBP) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.DeleteBP == nil {
			var val bool
			srv.DeleteBP = &val
		}

		*srv.DeleteBP, _ = strconv.ParseBool(*env.RAW_DeleteBP)
	}

	// Delete deaths override.
	if env.RAW_DeleteDeaths != nil && len(*env.RAW_DeleteDeaths) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.DeleteDeaths == nil {
			var val bool
			srv.DeleteDeaths = &val
		}

		*srv.DeleteDeaths, _ = strconv.ParseBool(*env.RAW_DeleteDeaths)
	}

	// Delete states override.
	if env.RAW_DeleteStates != nil && len(*env.RAW_DeleteStates) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.DeleteStates == nil {
			var val bool
			srv.DeleteStates = &val
		}

		*srv.DeleteStates, _ = strconv.ParseBool(*env.RAW_DeleteStates)
	}

	// Delete identities override.
	if env.RAW_DeleteIdentities != nil && len(*env.RAW_DeleteIdentities) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.DeleteIdentities == nil {
			var val bool
			srv.DeleteIdentities = &val
		}

		*srv.DeleteIdentities, _ = strconv.ParseBool(*env.RAW_DeleteIdentities)
	}
	// Delete tokens override.
	if env.RAW_DeleteTokens != nil && len(*env.RAW_DeleteTokens) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.DeleteTokens == nil {
			var val bool
			srv.DeleteTokens = &val
		}

		*srv.DeleteTokens, _ = strconv.ParseBool(*env.RAW_DeleteTokens)
	}

	// Merge delete files override.
	if env.RAW_DeleteFilesMerge != nil && len(*env.RAW_DeleteFilesMerge) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.DeleteFilesMerge == nil {
			var val bool
			srv.DeleteFilesMerge = &val
		}

		*srv.DeleteFilesMerge, _ = strconv.ParseBool(*env.RAW_DeleteFilesMerge)
	}

	// Additional file deletions
	if env.RAW_DeleteFiles != nil && len(*env.RAW_DeleteFiles) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.DeleteFiles == nil {
			var val []config.AdditionalFiles
			srv.DeleteFiles = &val
		}

		// Parse as string.
		data := *env.RAW_DeleteFiles

		// Create structure for expected format.
		var delete_file []config.AdditionalFiles

		// Convert string to structure via JSON.
		err := json.Unmarshal([]byte(data), &delete_file)

		if err == nil {
			*srv.DeleteFiles = delete_file
		}
	}

	// Delete server files/data override.
	if env.RAW_DeleteSv != nil && len(*env.RAW_DeleteSv) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.DeleteSv == nil {
			var val bool
			srv.DeleteSv = &val
		}

		*srv.DeleteSv, _ = strconv.ParseBool(*env.RAW_DeleteSv)
	}

	// Change world info override.
	if env.RAW_ChangeWorldInfo != nil && len(*env.RAW_ChangeWorldInfo) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.ChangeWorldInfo == nil {
			var val bool
			srv.ChangeWorldInfo = &val
		}

		*srv.ChangeWorldInfo, _ = strconv.ParseBool(*env.RAW_ChangeWorldInfo)
	}

	// World info override.
	if env.RAW_WorldInfo != nil && len(*env.RAW_WorldInfo) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.WorldInfo == nil {
			var val []config.WorldInfo
			srv.WorldInfo = &val
		}

		s := *env.RAW_WorldInfo
		var tmp []config.WorldInfo

		// Try to parse as JSON, if fails, parse as string.
		err := json.Unmarshal([]byte(s), &tmp)

		if err == nil {
			*srv.WorldInfo = tmp
		}
	}

	// Change world info pick type override.
	if env.RAW_WorldInfoPickType != nil && len(*env.RAW_WorldInfoPickType) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.WorldInfoPickType == nil {
			var val int
			srv.WorldInfoPickType = &val
		}

		val, _ := strconv.ParseInt(*env.RAW_WorldInfoPickType, 10, 16)

		*srv.WorldInfoPickType = int(val)
	}

	// Change world info merge override.
	if env.RAW_WorldInfoMerge != nil && len(*env.RAW_WorldInfoMerge) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.WorldInfoMerge == nil {
			var val bool
			srv.WorldInfoMerge = &val
		}

		*srv.WorldInfoMerge, _ = strconv.ParseBool(*env.RAW_WorldInfoMerge)
	}

	// Change hostname override.
	if env.RAW_ChangeHostname != nil && len(*env.RAW_ChangeHostname) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.ChangeHostName == nil {
			var val bool
			srv.ChangeHostName = &val
		}

		*srv.ChangeHostName, _ = strconv.ParseBool(*env.RAW_ChangeHostname)
	}

	// Hostname override.
	if env.RAW_Hostname != nil && len(*env.RAW_Hostname) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.HostName == nil {
			var val string
			srv.HostName = &val
		}

		*srv.HostName = *env.RAW_Hostname
	}

	// Merge warnings override.
	if env.RAW_MergeWarnings != nil && len(*env.RAW_MergeWarnings) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.MergeWarnings == nil {
			var val bool
			srv.MergeWarnings = &val
		}

		*srv.MergeWarnings, _ = strconv.ParseBool(*env.RAW_MergeWarnings)
	}

	// Wipe first override.
	if env.RAW_WipeFirst != nil && len(*env.RAW_WipeFirst) > 0 {
		srv.WipeFirst, _ = strconv.ParseBool(*env.RAW_WipeFirst)
	}

	// Warning messages override (another special case).
	if env.RAW_WarningMessages != nil && len(*env.RAW_WarningMessages) > 0 {
		// Make sure we don't need to allocate memory.
		if srv.WarningMessages == nil {
			var val []config.WarningMessage
			srv.WarningMessages = &val
		}

		// Parse as string.
		data := *env.RAW_WarningMessages

		// Create structure for expected format.
		var warning_msg []config.WarningMessage

		// Convert string to structure via JSON.
		err := json.Unmarshal([]byte(data), &warning_msg)

		if err == nil {
			*srv.WarningMessages = warning_msg
		}
	}

	return err
}
