package wipe

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/misc"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

// Processes files. Deletes files depending on the type of wipe.
func ProcessFiles(data *Data, UUID string) bool {
	var files_to_delete []string

	// Make sure to URL encode the query string (directory path).
	dir := url.QueryEscape(data.PathToServerFiles)

	ep := "client/servers/" + UUID + "/files/list?directory=" + dir

	// We first need to retrieve the current variable.
	d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "GET", ep, nil)

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => nil.")
	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "List Files return data => "+d+".")

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not list files in directory ("+dir+"). Please enable debugging level 4 for body response including errors.")

		return false
	}

	if err != nil {
		fmt.Println(err)

		return false
	}

	// We want to parse the response with the startup response structure.
	var files_list pterodactyl.ListFilesResp

	// Convert to JSON.
	err = json.Unmarshal([]byte(d), &files_list)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Loop through all files.
	for _, file := range files_list.Data {
		// Make sure we're not dealing with a directory or link.
		if !file.Attributes.IsFile || file.Attributes.IsSymlink {
			continue
		}

		add_to_del := false

		// Check if we want to wipe the map/save files.
		if data.DeleteMap {
			if strings.Contains(file.Attributes.Name, ".sav") || strings.Contains(file.Attributes.Name, ".map") {
				// Append file to list of files to delete.
				add_to_del = true
			}
		}

		// Check if we want to wipe blue prints.
		if data.DeleteBP {
			if strings.Contains(file.Attributes.Name, "blueprints") {
				add_to_del = true
			}
		}

		// Check if we want to wipe deaths.
		if data.DeleteDeaths {
			if strings.Contains(file.Attributes.Name, "deaths") {
				add_to_del = true
			}
		}

		// Check if we want to wipe states.
		if data.DeleteStates {
			if strings.Contains(file.Attributes.Name, "states") {
				add_to_del = true
			}
		}

		// Check if we want to wipe identities.
		if data.DeleteIdentities {
			if strings.Contains(file.Attributes.Name, "identities") {
				add_to_del = true
			}
		}

		// Check if we want to wipe tokens.
		if data.DeleteTokens {
			if strings.Contains(file.Attributes.Name, "tokens") {
				add_to_del = true
			}
		}

		// Check if we want to wipe server files/data.
		if data.DeleteSv {
			if strings.Contains(file.Attributes.Name, "sv.files") {
				add_to_del = true
			}
		}

		// If we want to delete the file, add it to the delete list.
		if add_to_del {
			files_to_delete = append(files_to_delete, file.Attributes.Name)
		}
	}

	// Prepare to delete files.
	post_data := make(map[string]interface{})
	post_data["root"] = data.PathToServerFiles
	post_data["files"] = files_to_delete

	// Debug.
	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Deleting files => "+strings.Join(files_to_delete, ", ")+" (directory = "+data.PathToServerFiles+").")

	ep = "client/servers/" + UUID + "/files/delete"

	// We first need to retrieve the current variable.
	d, _, err = pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "POST", ep, post_data)

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => "+misc.CreateKeyPairs(post_data)+".")
	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Delete Files return data => "+d+".")

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not delete files in directory ("+dir+"). Please enable debugging level 4 for body response including errors.")

		return false
	}

	if err != nil {
		fmt.Println(err)

		return false
	}

	return true
}
