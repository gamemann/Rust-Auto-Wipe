package format

import (
	"strconv"
	"strings"
	"time"
)

// Formats hostname for our needs.
func FormatString(str *string) {
	// Gain some information about the time.
	month := time.Now().Month()
	//week_day := time.Now().Weekday()
	day := time.Now().Day()
	hour := time.Now().Hour()
	min := time.Now().Minute()
	sec := time.Now().Second()

	*str = strings.Replace(*str, "{month}", strconv.Itoa(int(month)), -1)
	*str = strings.Replace(*str, "{day}", strconv.Itoa(day), -1)
	*str = strings.Replace(*str, "{hour}", strconv.Itoa(hour), -1)
	*str = strings.Replace(*str, "{minute}", strconv.Itoa(min), -1)
	*str = strings.Replace(*str, "{second}", strconv.Itoa(sec), -1)
}
