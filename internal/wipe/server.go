package wipe

import (
	"encoding/json"
	"errors"

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

	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Running State return data => "+d+".")

	if err != nil {
		return running, err
	}

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not get server's running state. Please enable debugging level 4 for body response including errors.")

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

	d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "POST", "client/servers/"+UUID+"/power", post_data)

	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Power Command return data => "+d+".")

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not send power command. Please enable debugging level 4 for body response including errors.")

		return errors.New("Could not send power command.")
	}

	return err
}
