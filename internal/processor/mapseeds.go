package processor

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gamemann/Rust-Auto-Wipe/internal/pterodactyl"
)

// Processes seeds and determines the next seed. Should occur before wipe.
func (wipedata *WipeData) ProcessSeeds(UUID string) bool {
	// We first need to retrieve the current variable.
	d, _, err := pterodactyl.SendAPIRequest(wipedata.APIURL, wipedata.APIToken, "GET", "client/servers/"+UUID+"/startup", nil)

	if err != nil {
		fmt.Println(err)

		return false
	}

	var EnvVars pterodactyl.StartupResp

	err = json.Unmarshal([]byte(d), &EnvVars)

	if err != nil {
		fmt.Println(err)

		return false
	}
	return true
}

// Gets the next seed in the array.
func (wipedata *WipeData) GetNextSeed(seeds []int, curseed int, picktype int) int {
	seed := -1

	if picktype == 1 {
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
