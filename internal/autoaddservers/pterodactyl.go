package autoaddservers

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

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
				StartupCommand string `json:"startup_command"`
				Image          string `json:"image"`
				Installed      bool   `json:"installed"`
				Environment    struct {
					ServerJarfile   string `json:"SERVER_JARFILE"`
					VanillaVersion  string `json:"VANILLA_VERSION"`
					Startup         string `json:"STARTUP"`
					PServerLocation string `json:"P_SERVER_LOCATION"`
					PServerUUID     string `json:"P_SERVER_UUID"`
				} `json:"environment"`
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
			var srv config.Server

			// Split UUID by -.
			uuid_split := strings.Split(v.Attributes.UUID, "-")

			// Assign short UUID.
			srv.UUID = uuid_split[0]

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
