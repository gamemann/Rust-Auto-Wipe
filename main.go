package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/internal/wipe"
)

func srv_handler(cfg *config.Config, srv *config.Server, idx int) {
	// We need to retrieve the wipe data information first.
	var data wipe.Data

	// For readability, define pointer to last week number here.
	last_day_num := &data.InternalData.LastDayNum
	last_month_num := &data.InternalData.LastMonthNum

	wipe.ProcessData(&data, cfg, srv)

	skip_next := false

	// Create a repeating loop until the two signals are called in the main function.
	for true {
		// We need to create our own time management.
		month := time.Now().Month()
		day := time.Now().Weekday()
		hour := time.Now().Hour()
		min := time.Now().Minute()

		new_month := false
		new_week := false

		do_wipe := false

		// Check if it's a new month (starts from 1 - 12).
		if *last_month_num != int(month) {
			new_month = true
		}

		if new_month && data.WipeMonthly {
			do_wipe = true
		}

		// First check special options.
		if new_week && data.WipeBiweekly {

		}

		// Otherwise, assume weekly. Check if we need to wipe.
		if uint8(day) == data.WipeDay && uint8(hour) == data.WipeHour && uint8(min) == data.WipeHour {
			// Check if we're doing bi-weekly.
			if data.WipeBiweekly {
				// Flip a switch.
				if !skip_next {
					do_wipe = true
					skip_next = true
				} else {
					skip_next = false
				}
			} else {
				// Otherwise, return true since we're assuming weekly.
				do_wipe = true
			}
		}

		// Check if we need to wipe.
		if do_wipe {
			// First, we can process seeds and host name.
			wipe.ProcessSeeds(&data, srv.UUID)
			wipe.ProcessHostName(&data, srv.UUID)

		}

		// Update last values.
		*last_day_num = int(day)
		*last_month_num = int(month)

		time.Sleep(time.Duration(time.Second))
	}
}

func main() {
	// Look for 'cfg' flag in command line arguments (default path: /etc/raw/raw.conf).
	configFile := flag.String("cfg", "/etc/raw/raw.conf", "The path to the Rust Auto Wipe config file.")
	flag.Parse()

	// Create config struct.
	cfg := config.Config{}

	// Set config defaults.
	cfg.SetDefaults()

	// Attempt to read config.
	cfg.LoadConfig(*configFile)

	// Loop through each server and execute Go routine.
	for i, srv := range cfg.Servers {
		go srv_handler(&cfg, &srv, i)
	}

	// Signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
}
