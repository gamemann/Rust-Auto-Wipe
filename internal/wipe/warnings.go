package wipe

import "github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"

func SendMessage(data *Data, UUID string, message string) error {
	var err error

	// Now convert to proper POST data.
	var post_data map[string]interface{}
	post_data["command"] = "say " + message

	// Send API request to update host name variable.
	_, _, err = pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "POST", "client/servers/"+UUID+"/POST", post_data)

	return err
}
