package autoaddservers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

type WarningMessageJSON struct {
	WarningTime uint   `json:"warningtime"`
	Message     string `json:"message"`
}

type WarningMessageOverride struct {
	Data []WarningMessageJSON `json:"data"`
}

type ServerListResp struct {
	Object string `json:"object"`
	Data   []struct {
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
				StartupCommand string `json:"startup_command"`
				Image          string `json:"image"`
				Installed      int    `json:"installed"`
				Environment    struct {
					WorldSeed *string `json:"WORLD_SEED"`
					HostName  *string `json:"HOSTNAME"`

					RAW_Enabled           *string `json:"RAW_ENABLED"`
					RAW_PathToServerFiles *string `json:"RAW_PATHTOSERVERFILES"`
					RAW_Timezone          *string `json:"RAW_TIMEZONE"`
					RAW_CronMerge         *string `json:"RAW_CRONMERGE"`
					RAW_CronStr           *string `json:"RAW_CRONSTR"`
					RAW_DeleteMap         *string `json:"RAW_DELETEMAP"`
					RAW_DeleteBP          *string `json:"RAW_DELETEBP"`
					RAW_DeleteDeaths      *string `json:"RAW_DELETEDEATHS"`
					RAW_DeleteStates      *string `json:"RAW_DELETESTATES"`
					RAW_DeleteIdentities  *string `json:"RAW_DELETEIDENTITIES"`
					RAW_DeleteTokens      *string `json:"RAW_DELETETOKENS"`
					RAW_DeleteSv          *string `json:"RAW_DELETESV"`
					RAW_ChangeMapSeeds    *string `json:"RAW_CHANGEMAPSEEDS"`
					RAW_MapSeeds          *string `json:"RAW_MAPSEEDS"`
					RAW_MapSeedsPickType  *string `json:"RAW_MAPSEEDSPICKTYPE"`
					RAW_MapSeedsMerge     *string `json:"RAW_MAPSEEDSMERGE"`
					RAW_ChangeHostname    *string `json:"RAW_CHANGEHOSTNAME"`
					RAW_Hostname          *string `json:"RAW_HOSTNAME"`
					RAW_MergeWarnings     *string `json:"RAW_MERGEWARNINGS"`
					RAW_WarningMessages   *string `json:"RAW_WARNINGMESSAGES"`
					RAW_WipeFirst         *string `json:"RAW_WIPEFIRST"`
				} `json:"environment"`
			} `json:"container"`
			UpdatedAt time.Time `json:"updated_at"`
			CreatedAt time.Time `json:"created_at"`
		} `json:"attributes"`
	} `json:"data"`
	Meta struct {
		Pagination struct {
			Total       int `json:"total"`
			Count       int `json:"count"`
			PerPage     int `json:"per_page"`
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
			Links       struct {
			} `json:"links"`
		} `json:"pagination"`
	} `json:"meta"`
}

func AddServers(cfg *config.Config) error {
	var err error

	// Page number.
	p := 1

	// Retrieve list of all servers from Pterodactyl application API.
	for true {
		ep := "application/servers?p=" + strconv.Itoa(p)

		d, _, err := pterodactyl.SendAPIRequest(cfg.APIURL, cfg.AppToken, "GET", ep, nil)

		debug.SendDebugMsg("AUTOADD", cfg.DebugLevel, 3, "Sending request. Request => "+cfg.APIURL+"api/"+ep+". Post data => nil.")
		debug.SendDebugMsg("AUTOADD", cfg.DebugLevel, 4, "Update Variable return data => "+d+".")

		if err != nil {
			break
		}

		// Convert JSON to structure.
		var server_list ServerListResp

		err = json.Unmarshal([]byte(d), &server_list)

		if err != nil {
			break
		}

		// Now loop through each data object (server).
		for _, v := range server_list.Data {
			// We must make sure the Rust environmental variables are valid if we're going to add said server.
			env := &v.Attributes.Container.Environment

			// If WORLD_SEED doesn't exist (empty field), don't add server.
			if env.WorldSeed == nil {
				continue
			}

			// If HOSTNAME doesn't exist (empty field), don't add server.
			if env.HostName == nil {
				continue
			}

			var srv config.Server

			// Enable by default.
			srv.Enabled = true

			// Split UUID by -.
			uuid_split := strings.Split(v.Attributes.UUID, "-")

			// Assign short UUID.
			srv.UUID = uuid_split[0]

			/* STARTING ENV OVERRIDES (using the opposite method as above) */
			// Enabled override.
			if env.RAW_Enabled != nil {
				srv.Enabled, _ = strconv.ParseBool(*env.RAW_Enabled)
			}

			// Path to server files override.
			if env.RAW_PathToServerFiles != nil {
				// Make sure we don't need to allocate memory.
				if srv.PathToServerFiles == nil {
					var val string
					srv.PathToServerFiles = &val
				}

				*srv.PathToServerFiles = *env.RAW_PathToServerFiles
			}

			// Timezone override.
			if env.RAW_Timezone != nil {
				// Make sure we don't need to allocate memory.
				if srv.Timezone == nil {
					var val string
					srv.Timezone = &val
				}

				*srv.Timezone = *env.RAW_Timezone
			}

			// Cron merge override.
			if env.RAW_CronMerge != nil {
				// Make sure we don't need to allocate memory.
				if srv.CronMerge == nil {
					var val bool
					srv.CronMerge = &val
				}

				*srv.CronMerge, _ = strconv.ParseBool(*env.RAW_CronMerge)
			}

			// Cron string override.
			if env.RAW_CronStr != nil {
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
			if env.RAW_DeleteMap != nil {
				// Make sure we don't need to allocate memory.
				if srv.DeleteMap == nil {
					var val bool
					srv.DeleteMap = &val
				}

				*srv.DeleteMap, _ = strconv.ParseBool(*env.RAW_DeleteMap)
			}

			// Delete blueprints override.
			if env.RAW_DeleteBP != nil {
				// Make sure we don't need to allocate memory.
				if srv.DeleteBP == nil {
					var val bool
					srv.DeleteBP = &val
				}

				*srv.DeleteBP, _ = strconv.ParseBool(*env.RAW_DeleteBP)
			}

			// Delete deaths override.
			if env.RAW_DeleteDeaths != nil {
				// Make sure we don't need to allocate memory.
				if srv.DeleteDeaths == nil {
					var val bool
					srv.DeleteDeaths = &val
				}

				*srv.DeleteDeaths, _ = strconv.ParseBool(*env.RAW_DeleteDeaths)
			}

			// Delete states override.
			if env.RAW_DeleteStates != nil {
				// Make sure we don't need to allocate memory.
				if srv.DeleteStates == nil {
					var val bool
					srv.DeleteStates = &val
				}

				*srv.DeleteStates, _ = strconv.ParseBool(*env.RAW_DeleteStates)
			}

			// Delete identities override.
			if env.RAW_DeleteIdentities != nil {
				// Make sure we don't need to allocate memory.
				if srv.DeleteIdentities == nil {
					var val bool
					srv.DeleteIdentities = &val
				}

				*srv.DeleteIdentities, _ = strconv.ParseBool(*env.RAW_DeleteIdentities)
			}
			// Delete tokens override.
			if env.RAW_DeleteTokens != nil {
				// Make sure we don't need to allocate memory.
				if srv.DeleteTokens == nil {
					var val bool
					srv.DeleteTokens = &val
				}

				*srv.DeleteTokens, _ = strconv.ParseBool(*env.RAW_DeleteTokens)
			}

			// Delete server files/data override.
			if env.RAW_DeleteSv != nil {
				// Make sure we don't need to allocate memory.
				if srv.DeleteSv == nil {
					var val bool
					srv.DeleteSv = &val
				}

				*srv.DeleteMap, _ = strconv.ParseBool(*env.RAW_DeleteSv)
			}

			// Change map seeds override.
			if env.RAW_ChangeMapSeeds != nil {
				// Make sure we don't need to allocate memory.
				if srv.ChangeMapSeeds == nil {
					var val bool
					srv.ChangeMapSeeds = &val
				}

				*srv.ChangeMapSeeds, _ = strconv.ParseBool(*env.RAW_ChangeMapSeeds)
			}

			// Map seeds override (this is a special case).
			if env.RAW_MapSeeds != nil {
				// Make sure we don't need to allocate memory.
				if srv.MapSeeds == nil {
					var val interface{}
					srv.MapSeeds = &val
				}

				s := *env.RAW_MapSeeds
				var tmp []int

				// Try to parse as JSON, if fails, parse as string.
				err := json.Unmarshal([]byte(s), &tmp)

				if err != nil {
					new_val, err := strconv.Atoi(*env.RAW_MapSeeds)

					if err == nil {
						*srv.MapSeeds = new_val
					}

				} else {
					*srv.MapSeeds = tmp
				}
			}

			// Change map seeds pick type override.
			if env.RAW_MapSeedsPickType != nil {
				// Make sure we don't need to allocate memory.
				if srv.MapSeedsPickType == nil {
					var val int
					srv.MapSeedsPickType = &val
				}

				val, _ := strconv.ParseInt(*env.RAW_MapSeedsPickType, 10, 16)

				*srv.MapSeedsPickType = int(val)
			}

			// Change map seeds merge override.
			if env.RAW_MapSeedsMerge != nil {
				// Make sure we don't need to allocate memory.
				if srv.MapSeedsMerge == nil {
					var val bool
					srv.MapSeedsMerge = &val
				}

				*srv.MapSeedsMerge, _ = strconv.ParseBool(*env.RAW_MapSeedsMerge)
			}

			// Change hostname override.
			if env.RAW_ChangeHostname != nil {
				// Make sure we don't need to allocate memory.
				if srv.ChangeHostName == nil {
					var val bool
					srv.ChangeHostName = &val
				}

				*srv.ChangeHostName, _ = strconv.ParseBool(*env.RAW_ChangeHostname)
			}

			// Hostname override.
			if env.RAW_Hostname != nil {
				// Make sure we don't need to allocate memory.
				if srv.HostName == nil {
					var val string
					srv.HostName = &val
				}

				*srv.HostName = *env.RAW_Hostname
			}

			// Merge warnings override.
			if env.RAW_MergeWarnings != nil {
				// Make sure we don't need to allocate memory.
				if srv.MergeWarnings == nil {
					var val bool
					srv.MergeWarnings = &val
				}

				*srv.MergeWarnings, _ = strconv.ParseBool(*env.RAW_MergeWarnings)
			}

			// Wipe first override.
			if env.RAW_WipeFirst != nil {
				srv.WipeFirst, _ = strconv.ParseBool(*env.RAW_WipeFirst)
			}

			// Warning messages override (another special case).
			if env.RAW_WarningMessages != nil {
				// Make sure we don't need to allocate memory.
				if srv.WarningMessages == nil {
					var val []config.WarningMessage
					srv.WarningMessages = &val
				}

				// Parse as string.
				data := *env.RAW_WarningMessages

				// Create structure for expected format.
				var warning_msg WarningMessageOverride

				// Convert string to structure via JSON.
				err := json.Unmarshal([]byte(data), &warning_msg)

				if err == nil {
					// Loop through entries and append to warning messages slice.
					for _, j := range warning_msg.Data {
						fmt.Println("AHHHHHH")
						var warning_msg_cfg config.WarningMessage
						warning_msg_cfg.WarningTime = j.WarningTime
						warning_msg_cfg.Message = j.Message

						*srv.WarningMessages = append(*srv.WarningMessages, warning_msg_cfg)
					}
				}
			}

			// Append to CFG server slice.
			cfg.Servers = append(cfg.Servers, srv)
		}

		// Check if we can exit now.
		if server_list.Meta.Pagination.CurrentPage >= server_list.Meta.Pagination.TotalPages {
			break
		}

		// Increment page number.
		p++
	}

	return err
}
