package wipe

import (
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

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

func SendPowerCommand(data *Data, UUID string, cmd string) error {
	var post_data map[string]interface{}
	post_data["signal"] = cmd

	_, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "POST", "client/servers/"+UUID+"/power", post_data)

	return err
}
