package wipe

import (
	"fmt"

	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/format"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

// Sets the next host name to use.
func ProcessHostName(data *Data, UUID string, month int, day int, week_day int) bool {
	hostname := data.HostName

	// Format hostname.
	format.FormatString(&hostname, month, day, week_day, 0, 0, 0)

	// Now convert to proper POST data.
	post_data := make(map[string]interface{})
	post_data["key"] = "HOSTNAME"
	post_data["value"] = hostname

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Setting hostname => \""+hostname+"\".")

	// Send API request to update host name variable.
	_, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "PUT", "client/servers/"+UUID+"/variable", post_data)

	if err != nil {
		fmt.Println(err)

		return false
	}

	return true
}
