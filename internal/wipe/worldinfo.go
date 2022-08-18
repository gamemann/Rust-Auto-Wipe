package wipe

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/misc"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/pterodactyl"
)

// Processes world info such as map, size, and seed.
func ProcessWorldInfo(data *Data, UUID string) bool {
	cur_world := &data.WorldInfo[data.InternalData.LatestWorld]

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "[WI] Current map => \""+*cur_world.Map+"\".")
	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "[WI] Current size => \""+strconv.Itoa(*cur_world.WorldSize)+"\".")
	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "[WI] Current seed => \""+strconv.Itoa(*cur_world.WorldSeed)+"\".")

	// Now get the next world using the GetNextWorld() method.
	next_world := GetNextWorld(data)

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "[WI] Next map => \""+*next_world.Map+"\".")
	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "[WI] Next size => \""+strconv.Itoa(*next_world.WorldSize)+"\".")
	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "[WI] Next seed => \""+strconv.Itoa(*next_world.WorldSeed)+"\".")

	// First, set world map.
	post_data := make(map[string]interface{})
	post_data["key"] = "LEVEL"
	post_data["value"] = *next_world.Map

	ep := "client/servers/" + UUID + "/startup/variable"

	// Send API request.
	d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "PUT", ep, post_data)

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => "+misc.CreateKeyPairs(post_data)+".")
	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Update Variable (LEVEL) return data => "+d+".")

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not update startup LEVEL variable. Please enable debugging level 4 for body response including errors.")

		return false
	}

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Now, set world size.
	post_data = make(map[string]interface{})
	post_data["key"] = "WORLD_SIZE"
	post_data["value"] = strconv.Itoa(*next_world.WorldSize)

	ep = "client/servers/" + UUID + "/startup/variable"

	// Send API request.
	d, _, err = pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "PUT", ep, post_data)

	debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => "+misc.CreateKeyPairs(post_data)+".")
	debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Update Variable (WORLD_SIZE) return data => "+d+".")

	if pterodactyl.IsError(d) {
		debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not update startup WORLD_SIZE variable. Please enable debugging level 4 for body response including errors.")

		return false
	}

	if err != nil {
		fmt.Println(err)

		return false
	}

	// Finally, set the world seed if needed.
	if *next_world.WorldSeed > 0 {
		post_data := make(map[string]interface{})
		post_data["key"] = "WORLD_SEED"
		post_data["value"] = strconv.Itoa(*next_world.WorldSeed)

		ep := "client/servers/" + UUID + "/startup/variable"

		// Send API request.
		d, _, err := pterodactyl.SendAPIRequest(data.APIURL, data.APIToken, "PUT", ep, post_data)

		debug.SendDebugMsg(UUID, data.DebugLevel, 3, "Sending request. Request => "+data.APIURL+"api/"+ep+". Post data => "+misc.CreateKeyPairs(post_data)+".")
		debug.SendDebugMsg(UUID, data.DebugLevel, 4, "Update Variable (WORLD_SEED) return data => "+d+".")

		if pterodactyl.IsError(d) {
			debug.SendDebugMsg(UUID, data.DebugLevel, 0, "Could not update startup WORLD_SEED variable. Please enable debugging level 4 for body response including errors.")

			return false
		}

		if err != nil {
			fmt.Println(err)

			return false
		}
	}

	return true
}

// Gets the world from the array.
func GetNextWorld(data *Data) config.WorldInfo {
	// Make new variables for better looking code.
	world := config.WorldInfo{}
	worlds := data.WorldInfo
	pick_type := data.WorldInfoPickType

	// Check pick type.
	if pick_type == 1 {
		next_idx := &data.InternalData.LatestWorld

		world = data.WorldInfo[*next_idx]

		// Increment current (technically next) world.
		*next_idx++

		// Make sure we don't need to reset.
		if *next_idx >= uint(len(data.WorldInfo)) {
			*next_idx = 0
		}

	} else {
		rand.Seed(time.Now().UnixNano())
		min := 0
		max := len(worlds)
		world = data.WorldInfo[rand.Intn(max-min+1)+min]
	}

	return world
}
