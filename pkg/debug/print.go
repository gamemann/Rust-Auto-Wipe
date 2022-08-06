package debug

import (
	"fmt"
	"strconv"
)

func SendDebugMsg(UUID string, debug_level int, required_level int, message string) {
	if debug_level >= required_level {
		pre := "ERR"

		if required_level > 0 {
			pre = strconv.Itoa(required_level)
		}
		fmt.Println("[" + pre + "] " + UUID + " :: " + message)
	}
}
