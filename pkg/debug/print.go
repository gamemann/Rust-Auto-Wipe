package debug

import (
	"fmt"
	"strconv"
	"time"
)

func SendDebugMsg(UUID string, debug_level int, required_level int, message string) {
	if debug_level >= required_level {
		now := time.Now().Local()

		pre := "ERR"

		if required_level > 0 {
			pre = strconv.Itoa(required_level)
		}
		fmt.Println("[" + pre + "][" + now.Format("1-2 03:04:05 PM") + "] " + UUID + " :: " + message)
	}
}
