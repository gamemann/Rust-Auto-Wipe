package debug

import (
	"fmt"
	"strconv"
)

func SendDebugMsg(UUID string, debug_level int, required_level int, message string) {
	if debug_level >= required_level {
		fmt.Println("[" + strconv.Itoa(required_level) + "] " + UUID + " :: " + message)
	}
}
