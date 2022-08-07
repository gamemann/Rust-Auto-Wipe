package format

import (
	"strconv"
	"strings"
	"time"
)

// Formats hostname for our needs.
func FormatString(str *string, secs_left int) {
	*str = strings.Replace(*str, "{seconds_left}", strconv.Itoa(secs_left), -1)

	// Gain some information about the time.
	tz_one := time.Now().Format("MST")
	tz_two := time.Now().Format("-0700")
	tz_three := time.Now().Format("-07")

	month_str_short := time.Now().Format("Jan")
	month_str_long := time.Now().Format("January")

	week_day_str_short := time.Now().Format("Mon")
	week_day_str_long := time.Now().Format("Monday")

	year_one := time.Now().Format("06")
	year_two := time.Now().Format("2006")

	month_one := time.Now().Format("01")
	month_two := time.Now().Format("1")
	month_three := time.Now().Format("_1")

	day_one := time.Now().Format("02")
	day_two := time.Now().Format("2")
	day_three := time.Now().Format("_2")

	hour_one := time.Now().Format("03")
	hour_two := time.Now().Format("3")
	hour_three := time.Now().Format("15")

	min_one := time.Now().Format("04")
	min_two := time.Now().Format("4")

	sec_one := time.Now().Format("05")
	sec_two := time.Now().Format("5")

	mark_one := time.Now().Format("PM")
	mark_two := time.Now().Format("pm")

	*str = strings.Replace(*str, "{tz_one}", tz_one, -1)
	*str = strings.Replace(*str, "{tz_two}", tz_two, -1)
	*str = strings.Replace(*str, "{tz_three}", tz_three, -1)

	*str = strings.Replace(*str, "{month_str_short}", month_str_short, -1)
	*str = strings.Replace(*str, "{month_str_long}", month_str_long, -1)

	*str = strings.Replace(*str, "{week_day_str_short}", week_day_str_short, -1)
	*str = strings.Replace(*str, "{week_day_str_long}", week_day_str_long, -1)

	*str = strings.Replace(*str, "{year_one}", year_one, -1)
	*str = strings.Replace(*str, "{year_two}", year_two, -1)

	*str = strings.Replace(*str, "{month_one}", month_one, -1)
	*str = strings.Replace(*str, "{month_two}", month_two, -1)
	*str = strings.Replace(*str, "{month_three}", month_three, -1)

	*str = strings.Replace(*str, "{day_one}", day_one, -1)
	*str = strings.Replace(*str, "{day_two}", day_two, -1)
	*str = strings.Replace(*str, "{day_three}", day_three, -1)

	*str = strings.Replace(*str, "{hour_one}", hour_one, -1)
	*str = strings.Replace(*str, "{hour_two}", hour_two, -1)
	*str = strings.Replace(*str, "{hour_three}", hour_three, -1)

	*str = strings.Replace(*str, "{min_one}", min_one, -1)
	*str = strings.Replace(*str, "{min_two}", min_two, -1)

	*str = strings.Replace(*str, "{sec_one}", sec_one, -1)
	*str = strings.Replace(*str, "{sec_two}", sec_two, -1)

	*str = strings.Replace(*str, "{mark_one}", mark_one, -1)
	*str = strings.Replace(*str, "{mark_two}", mark_two, -1)
}
