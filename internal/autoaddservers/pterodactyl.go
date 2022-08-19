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

type RawEnv struct {
	WorldSeed  *string `json:"WORLD_SEED"`
	HostName   *string `json:"HOSTNAME"`
	ServerIP   *string `json:"SERVER_IP"`
	ServerPort *string `json:"SERVER_PORT"`

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
	RAW_DeleteFilesMerge  *string `json:"RAW_DELETEFILESMERGE"`
	RAW_DeleteFiles       *string `json:"RAW_DELETEFILES"`
	RAW_DeleteSv          *string `json:"RAW_DELETESV"`
	RAW_ChangeWorldInfo   *string `json:"RAW_CHANGEWORLDINFO"`
	RAW_WorldInfo         *string `json:"RAW_WORLDINFO"`
	RAW_WorldInfoPickType *string `json:"RAW_WORLDINFOPICKTYPE"`
	RAW_WorldInfoMerge    *string `json:"RAW_WORLDINFOMERGE"`
	RAW_ChangeHostname    *string `json:"RAW_CHANGEHOSTNAME"`
	RAW_Hostname          *string `json:"RAW_HOSTNAME"`
	RAW_MergeWarnings     *string `json:"RAW_MERGEWARNINGS"`
	RAW_WarningMessages   *string `json:"RAW_WARNINGMESSAGES"`
	RAW_WipeFirst         *string `json:"RAW_WIPEFIRST"`
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
				Environment    RawEnv `json:"environment"`
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
			var uuid_split []string
			var srv config.Server

			// We must make sure the Rust environmental variables are valid if we're going to add said server.
			env := &v.Attributes.Container.Environment

			ip := ""

			if env.ServerIP != nil {
				ip = *env.ServerIP
			}

			port := ""

			if env.ServerPort != nil {
				port = *env.ServerPort
			}

			// Loop through all current servers and make sure we update the ID if necessary.
			if !cfg.AutoAddServers {
				goto updateid
			}

			// If WORLD_SEED doesn't exist (empty field), don't add server.
			if env.WorldSeed == nil {
				continue
			}

			// If HOSTNAME doesn't exist (empty field), don't add server.
			if env.HostName == nil {
				continue
			}

			// Enable by default.
			srv.Enabled = true

			// Split UUID by -.
			uuid_split = strings.Split(v.Attributes.UUID, "-")

			// Assign short UUID.
			srv.UUID = uuid_split[0]

			// Append to CFG server slice.
			cfg.Servers = append(cfg.Servers, srv)

		updateid:
			// Get server information
			for i := 0; i < len(cfg.Servers); i++ {
				srv := &cfg.Servers[i]

				if srv.UUID == v.Attributes.Identifier {
					srv.ID = v.Attributes.ID
					srv.Name = v.Attributes.Name
					srv.LongID = v.Attributes.UUID
					srv.IP = ip
					srv.Port, err = strconv.Atoi(port)

					if err != nil {
						fmt.Println(err)
					}
				}
			}
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
