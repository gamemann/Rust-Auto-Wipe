package format

import (
	"strconv"
	"strings"
)

// Formats hostname for our needs.
func FormatString(str *string, month int, day int, week_day int, hour int, minute int, second int) {
	*str = strings.Replace(*str, "{month}", strconv.Itoa(month), -1)
	*str = strings.Replace(*str, "{day}", strconv.Itoa(day), -1)
	*str = strings.Replace(*str, "{hour}", strconv.Itoa(hour), -1)
	*str = strings.Replace(*str, "{minute}", strconv.Itoa(minute), -1)
	*str = strings.Replace(*str, "{second}", strconv.Itoa(second), -1)
}
