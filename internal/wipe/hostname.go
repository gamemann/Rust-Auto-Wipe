package wipe

import (
	"fmt"

	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/format"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/misc"
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

	ep := "client/servers/" + UUID + "/startup/variable"

	// Send API request to update host name variable.
	d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "PUT", ep, post_data)

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => "+misc.CreateKeyPairs(post_data)+".")
	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Update Variable return data => "+d+".")

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not update hostname. Please enable debugging level 4 for body response including errors.")

		return false
	}

	if err != nil {
		fmt.Println(err)

		return false
	}

	return true
}
