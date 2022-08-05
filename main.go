package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/internal/wipe"
)

func srv_handler(cfg *config.Config, srv *config.Server, idx int) {
	// We need to retrieve the wipe data information first.
	var data wipe.Data

	// For readability, define pointer to last week/day numbers here.
	last_day_num := &data.InternalData.LastDayNum
	last_month_num := &data.InternalData.LastMonthNum

	wipe.ProcessData(&data, cfg, srv)

	skip_next := true

	// Create a repeating loop until the two signals are called in the main function.
	for true {
		// We need to create our own time management.
		month := time.Now().Month()
		week_day := time.Now().Weekday()
		day := time.Now().Day()
		hour := time.Now().Hour()
		min := time.Now().Minute()

		new_month := false

		do_wipe := false

		// Check if it's a new month (starts from 1 - 12).
		if *last_month_num != int(month) {
			new_month = true
		}

		if new_month && data.WipeMonthly {
			do_wipe = true
		}

		// Otherwise, assume weekly. Check if we need to wipe.
		if uint8(week_day) == data.WipeDay && uint8(hour) == data.WipeHour && uint8(min) == data.WipeHour {
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
			// Process map seeds.
			if data.ChangeMapSeeds {
				wipe.ProcessSeeds(&data, srv.UUID)
			}

			// Process host name.
			if data.ChangeHostName {
				wipe.ProcessHostName(&data, srv.UUID, int(month), int(day), int(week_day))
			}

			// We should stop the server (To Do: Implement something to check if server is running and force kill if so).
			wipe.StopServer(&data, srv.UUID)

			// Wait until the server is confirmed stopped.
			i := 0

			for true {
				// Check if the server is running and when it is confirmed stop, break the loop.
				running, err := wipe.IsServerRunning(&data, srv.UUID)

				// Check for error. Otherwise, break if we're not running.
				if err != nil {
					fmt.Println(err)
				} else {
					if !running {
						break
					}
				}

				// Increment i.
				i++

				// Kill the server after 15 seconds.
				if i > 15 {
					wipe.KillServer(&data, srv.UUID)
				}

				// Give up after a minute.
				if i > 60 {
					break
				}

				// Sleep every second to avoid unnecessary CPU cycles.
				time.Sleep(time.Duration(time.Second))
			}

			// Process and delete files.
			wipe.ProcessFiles(&data, srv.UUID)

			// Start server back up.
			wipe.StartServer(&data, srv.UUID)
		}

		// Update last values.
		*last_day_num = int(week_day)
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
	err := cfg.LoadConfig(*configFile)

	// If we have no config, create the file with the defaults.
	if err != nil {
		// If there's an error and it contains "no such file", try to create the file with defaults.
		if strings.Contains(err.Error(), "no such file") {
			err = cfg.WriteDefaultsToFile(*configFile)

			if err != nil {
				fmt.Println("Failed to open config file and cannot create file.")
				fmt.Println(err)

				return
			}
		}

		fmt.Println("WARNING - No config file found. Created config file at " + *configFile + " with defaults.")
	}

	if len(cfg.Servers) < 1 {
		fmt.Println("[ERR] No servers found.")

		return
	}

	// Loop through each server and execute Go routine.
	for i, srv := range cfg.Servers {
		go srv_handler(&cfg, &srv, i)
	}

	// Signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
}
