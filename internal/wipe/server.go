package wipe

import (
	"encoding/json"

	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

type ServerResources struct {
	Object     string `json:"object"`
	Attributes struct {
		CurrentState string `json:"current_state"`
		IsSuspended  bool   `json:"is_suspended"`
		Resources    struct {
			MemoryBytes    int `json:"memory_bytes"`
			CPUAbsolute    int `json:"cpu_absolute"`
			DiskBytes      int `json:"disk_bytes"`
			NetworkRxBytes int `json:"network_rx_bytes"`
			NetworkTxBytes int `json:"network_tx_bytes"`
		} `json:"resources"`
	} `json:"attributes"`
}

func StartServer(data *Data, UUID string) error {

	err := SendPowerCommand(data, UUID, "start")

	return err
}

func StopServer(data *Data, UUID string) error {

	err := SendPowerCommand(data, UUID, "stop")

	return err
}

func KillServer(data *Data, UUID string) error {

	err := SendPowerCommand(data, UUID, "kill")

	return err
}

func IsServerRunning(data *Data, UUID string) (bool, error) {
	var running bool = false
	var err error

	d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "GET", "client/servers/"+UUID+"/resources", nil)

	if err != nil {
		return running, err
	}

	var resources ServerResources

	err = json.Unmarshal([]byte(d), &resources)

	if err != nil {
		return running, err
	}

	if resources.Attributes.CurrentState == "running" {
		running = true
	}

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Running state => "+resources.Attributes.CurrentState+".")

	return running, err
}

func SendPowerCommand(data *Data, UUID string, cmd string) error {
	post_data := make(map[string]interface{})
	post_data["signal"] = cmd

	_, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "POST", "client/servers/"+UUID+"/power", post_data)

	return err
}
