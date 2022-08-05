package processor

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

// Processes seeds and determines the next seed. Should occur before wipe.
func (wipedata *WipeData) ProcessSeeds(UUID string) bool {
	// We first need to retrieve the current variable.
	d, _, err := pterodactyl.SendAPIRequest(wipedata.APIURL, wipedata.APIToken, "GET", "client/servers/"+UUID+"/startup", nil)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// We want to parse the response with the startup response structure.
	var EnvVars pterodactyl.StartupResp

	// Convert to JSON.
	err = json.Unmarshal([]byte(d), &EnvVars)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Default seed = empty.
	cur_seed_str := ""

	// Loop through all startup variables and find the seed env value.
	for _, seed := range EnvVars.Data {
		if seed.Attributes.Env_Variable == "WORLD_SEED" {
			cur_seed_str = seed.Attributes.Srv_Value
		}
	}

	// Let's use ParseInt() for flexability.
	curseed64, err := strconv.ParseInt(cur_seed_str, 10, 16)

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Convert to integer type (To Do: Find a less sloppy way of doing this).
	cur_seed := int(curseed64)

	// Now get the next seed using the GetNextSeed() method.
	next_seed := wipedata.GetNextSeed(cur_seed)

	// Now convert to proper POST data.
	var post_data map[string]string
	post_data["key"] = "WORLD_SEED"
	post_data["value"] = strconv.Itoa(next_seed)

	// Send API request.
	d, _, err = pterodactyl.SendAPIRequest(wipedata.APIURL, wipedata.APIToken, "PUT", "client/servers/"+UUID+"/variable", post_data)

	return true
}

// Gets the next seed in the array.
func (wipedata *WipeData) GetNextSeed(curseed int) int {
	// Make new variables for better looking code.
	seed := -1
	seeds := wipedata.MapSeeds
	pick_type := wipedata.MapSeedPickType

	// Check pick type.
	if pick_type == 1 {
		// Loop through all seeds and get the next seed.
		for v, s := range seeds {
			if curseed == s {
				// If we're on the last seed, return 0 as the array item (starting item). Otherwise, return index + 1.
				if (len(seeds) - 1) == v {
					seed = 0
				} else {
					seed = v + 1
				}
			}
		}
	} else {
		seed = rand.Intn((len(seeds)-1)+1) + 0
	}

	return seed
}