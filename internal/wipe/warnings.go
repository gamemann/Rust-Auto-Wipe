package wipe

import (
	"errors"

	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

func SendMessage(data *Data, UUID string, message string) error {
	var err error

	// Now convert to proper POST data.
	post_data := make(map[string]interface{})
	post_data["command"] = "say " + message

	ep := "client/servers/" + UUID + "/command"

	// Send API request to update host name variable.
	d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "POST", ep, post_data)

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => nil.")
	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Command return data => "+d+".")

	if err != nil {
		return err
	}

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not send command to server (warning say). Please enable debugging level 4 for body response including errors.")

		return errors.New("Could not send command to server (warning say).")
	}

	return err
}
