package wipe

import (
	"encoding/json"
	"errors"

	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/misc"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

type ServerResources struct {
	Object     string `json:"object"`
	Attributes struct {
		CurrentState string `json:"current_state"`
		IsSuspended  bool   `json:"is_suspended"`
		Resources    struct {
			MemoryBytes    int64   `json:"memory_bytes"`
			CPUAbsolute    float64 `json:"cpu_absolute"`
			DiskBytes      int64   `json:"disk_bytes"`
			NetworkRxBytes int     `json:"network_rx_bytes"`
			NetworkTxBytes int     `json:"network_tx_bytes"`
			Uptime         int     `json:"uptime"`
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

func GetServerState(data *Data, UUID string) (string, error) {
	var state string = "stopped"
	var err error

	ep := "client/servers/" + UUID + "/resources"

	d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "GET", ep, nil)

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => nil.")
	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Running State return data => "+d+".")

	if err != nil {
		return state, err
	}

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not get server's running state. Please enable debugging level 4 for body response including errors.")

		return state, err
	}

	var resources ServerResources

	err = json.Unmarshal([]byte(d), &resources)

	if err != nil {
		return state, err
	}

	state = resources.Attributes.CurrentState

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Server state => "+state+".")

	return state, err
}

func IsServerRunning(data *Data, UUID string) (bool, error) {
	var running bool
	var err error

	state, err := GetServerState(data, UUID)

	if err != nil {
		return running, err
	}

	if state == "running" {
		running = true
	}

	return running, err
}

func SendPowerCommand(data *Data, UUID string, cmd string) error {
	post_data := make(map[string]interface{})
	post_data["signal"] = cmd

	ep := "client/servers/" + UUID + "/power"

	d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "POST", ep, post_data)

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => "+misc.CreateKeyPairs(post_data)+".")
	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Power Command return data => "+d+".")

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not send power command. Please enable debugging level 4 for body response including errors.")

		return errors.New("Could not send power command.")
	}

	return err
}
