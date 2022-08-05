package autoaddservers

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

type WarningMessageJSON struct {
	WarningTime uint   `json:"time"`
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
			ID          int    `json:"id"`
			ExternalID  string `json:"external_id"`
			UUID        string `json:"uuid"`
			Identifier  string `json:"identifier"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Suspended   bool   `json:"suspended"`
			Limits      struct {
				Memory  int         `json:"memory"`
				Swap    int         `json:"swap"`
				Disk    int         `json:"disk"`
				Io      int         `json:"io"`
				CPU     int         `json:"cpu"`
				Threads interface{} `json:"threads"`
			} `json:"limits"`
			FeatureLimits struct {
				Databases   int `json:"databases"`
				Allocations int `json:"allocations"`
				Backups     int `json:"backups"`
			} `json:"feature_limits"`
			User       int         `json:"user"`
			Node       int         `json:"node"`
			Allocation int         `json:"allocation"`
			Nest       int         `json:"nest"`
			Egg        int         `json:"egg"`
			Pack       interface{} `json:"pack"`
			Container  struct {
				StartupCommand string      `json:"startup_command"`
				Image          string      `json:"image"`
				Installed      bool        `json:"installed"`
				Environment    interface{} `json:"environment"`
			} `json:"container"`
			UpdatedAt     time.Time `json:"updated_at"`
			CreatedAt     time.Time `json:"created_at"`
			Relationships struct {
				Databases struct {
					Object string `json:"object"`
					Data   []struct {
						Object     string `json:"object"`
						Attributes struct {
							ID             int       `json:"id"`
							Server         int       `json:"server"`
							Host           int       `json:"host"`
							Database       string    `json:"database"`
							Username       string    `json:"username"`
							Remote         string    `json:"remote"`
							MaxConnections int       `json:"max_connections"`
							CreatedAt      time.Time `json:"created_at"`
							UpdatedAt      time.Time `json:"updated_at"`
						} `json:"attributes"`
					} `json:"data"`
				} `json:"databases"`
			} `json:"relationships"`
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
		d, _, err := pterodactyl.SendAPIRequest(cfg.APIURL, cfg.AppToken, "GET", "application/servers?p="+strconv.Itoa(p), nil)

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

			meta_val := reflect.ValueOf(*env).Elem()
			fld := meta_val.FieldByName("WORLD_SEED")

			// If WORLD_SEED doesn't exist (empty field), don't add server.
			if fld == (reflect.Value{}) {
				continue
			}

			fld = meta_val.FieldByName("HOSTNAME")

			// If HOSTNAME doesn't exist (empty field), don't add server.
			if fld == (reflect.Value{}) {
				continue
			}

			var srv config.Server

			// Split UUID by -.
			uuid_split := strings.Split(v.Attributes.UUID, "-")

			// Assign short UUID.
			srv.UUID = uuid_split[0]

			/* STARTING ENV OVERRIDES (using the opposite method as above) */
			// Enabled override.
			fld = meta_val.FieldByName("RAW_ENABLED")

			if fld != (reflect.Value{}) {
				srv.Enabled = reflect.Value.Bool(fld)
			}

			// Path to server files override.
			fld = meta_val.FieldByName("RAW_PATHTOSERVERFILES")

			if fld != (reflect.Value{}) {
				*srv.PathToServerFiles = reflect.Value.String(fld)
			}

			// Timezone override.
			fld = meta_val.FieldByName("RAW_TIMEZONE")

			if fld != (reflect.Value{}) {
				*srv.Timezone = reflect.Value.String(fld)
			}

			// Wipe time override.
			fld = meta_val.FieldByName("RAW_WIPETIME")

			if fld != (reflect.Value{}) {
				*srv.WipeTime = reflect.Value.String(fld)
			}

			// Wipe monthly override.
			fld = meta_val.FieldByName("RAW_MONTHLY")

			if fld != (reflect.Value{}) {
				*srv.WipeMonthly = reflect.Value.Bool(fld)
			}

			// Wipe bi-weekly override.
			fld = meta_val.FieldByName("RAW_BIWEEKLY")

			if fld != (reflect.Value{}) {
				*srv.WipeBiweekly = reflect.Value.Bool(fld)
			}

			// Delete map override.
			fld = meta_val.FieldByName("RAW_DELETEMAP")

			if fld != (reflect.Value{}) {
				*srv.DeleteMap = reflect.Value.Bool(fld)
			}

			// Delete blueprints override.
			fld = meta_val.FieldByName("RAW_DELETEBP")

			if fld != (reflect.Value{}) {
				*srv.DeleteBP = reflect.Value.Bool(fld)
			}

			// Delete deaths override.
			fld = meta_val.FieldByName("RAW_DELETEDEATHS")

			if fld != (reflect.Value{}) {
				*srv.DeleteDeaths = reflect.Value.Bool(fld)
			}

			// Delete states override.
			fld = meta_val.FieldByName("RAW_DELETESTATES")

			if fld != (reflect.Value{}) {
				*srv.DeleteStates = reflect.Value.Bool(fld)
			}

			// Delete identities override.
			fld = meta_val.FieldByName("RAW_DELETEIDENTITIES")

			if fld != (reflect.Value{}) {
				*srv.DeleteIdentities = reflect.Value.Bool(fld)
			}

			// Delete tokens override.
			fld = meta_val.FieldByName("RAW_DELETETOKENS")

			if fld != (reflect.Value{}) {
				*srv.DeleteTokens = reflect.Value.Bool(fld)
			}

			// Delete server files/data override.
			fld = meta_val.FieldByName("RAW_DELETESV")

			if fld != (reflect.Value{}) {
				*srv.DeleteSv = reflect.Value.Bool(fld)
			}

			// Change map seeds override.
			fld = meta_val.FieldByName("RAW_CHANGEMAPSEEDS")

			if fld != (reflect.Value{}) {
				*srv.ChangeMapSeeds = reflect.Value.Bool(fld)
			}

			// Map seeds override (this is a special case).
			fld = meta_val.FieldByName("RAW_MAPSEEDS")

			if fld != (reflect.Value{}) {
				// Parse as string and split by ",".
				seeds_str := reflect.Value.String(fld)
				seeds_split := strings.Split(seeds_str, ",")

				// Now loop through and insert into map seeds slice.
				for _, seed := range seeds_split {
					seed_num, err := strconv.Atoi(seed)

					if err != nil {
						continue
					}

					*srv.MapSeeds = append(*srv.MapSeeds, seed_num)
				}
			}

			// Change map seeds pick type override.
			fld = meta_val.FieldByName("RAW_MAPSEEDSPICKTYPE")

			if fld != (reflect.Value{}) {
				*srv.MapSeedsPickType = int(reflect.Value.Int(fld))
			}

			// Change map seeds pick type override.
			fld = meta_val.FieldByName("RAW_MAPSEEDSMERGE")

			if fld != (reflect.Value{}) {
				*srv.MapSeedsMerge = reflect.Value.Bool(fld)
			}

			// Change hostname override.
			fld = meta_val.FieldByName("RAW_CHANGEHOSTNAME")

			if fld != (reflect.Value{}) {
				*srv.ChangeHostName = reflect.Value.Bool(fld)
			}

			// Hostname override.
			fld = meta_val.FieldByName("RAW_HOSTNAME")

			if fld != (reflect.Value{}) {
				*srv.HostName = reflect.Value.String(fld)
			}

			// Merge warnings override.
			fld = meta_val.FieldByName("RAW_MERGEWARNINGS")

			if fld != (reflect.Value{}) {
				*srv.MergeWarnings = reflect.Value.Bool(fld)
			}

			// Warning messages override (another special case).
			fld = meta_val.FieldByName("RAW_MERGEWARNINGS")

			if fld != (reflect.Value{}) {
				// Parse as string.
				data := reflect.Value.String(fld)

				// Create structure for expected format.
				var warning_msg WarningMessageOverride

				// Convert string to structure via JSON.
				err := json.Unmarshal([]byte(data), &warning_msg)

				if err == nil {
					// Loop through entries and append to warning messages slice.
					for _, j := range warning_msg.Data {
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
