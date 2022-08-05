package wipe

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

// Sets the next host name to use.
func ProcessHostName(data *Data, UUID string, month int, day int, week_day int) bool {
	hostname := data.HostName

	// Format hostname.
	FormatHostname(&hostname, month, day, week_day)

	// Now convert to proper POST data.
	var post_data map[string]string
	post_data["key"] = "HOSTNAME"
	post_data["value"] = hostname

	// Send API request to update host name variable.
	_, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "PUT", "client/servers/"+UUID+"/variable", post_data)

	if err != nil {
		fmt.Println(err)

		return false
	}

	return true
}

// Formats hostname for our needs.
func FormatHostname(hostname *string, month int, day int, week_day int) {
	*hostname = strings.Replace(*hostname, "{month}", strconv.Itoa(month), -1)
	*hostname = strings.Replace(*hostname, "{day}", strconv.Itoa(day), -1)
}
