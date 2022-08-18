package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/internal/autoaddservers"
	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/internal/wipe"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/format"
	cron "github.com/robfig/cron/v3"
)

const HELP_MENU = "Help Options\n\t-cfg= --cfg -cfg <path> > Path to config file override.\n\t-l --list > Print out full config.\n\t-v --version > Print out version and exit.\n\t-h --help > Display help menu.\n\n"
const VERSION = "1.0.0"

func wipe_server(cfg *config.Config, srv *config.Server, data *wipe.Data) {
	debug.SendDebugMsg(srv.UUID, data.DebugLevel, 1, "Wiping server...")

	// Process world info.
	if data.ChangeWorldInfo {
		debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Processing world info...")

		wipe.ProcessWorldInfo(data, srv.UUID)
	}

	// Process host name.
	if data.ChangeHostName {
		debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Processing hostname...")

		wipe.ProcessHostName(data, srv.UUID)
	}

	debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Stopping server...")

	// We should stop the server
	wipe.StopServer(data, srv.UUID)

	// Wait until the server is confirmed stopped.
	i := 0

	for true {
		// Check if the server is running and when it is confirmed stop, break the loop.
		state, err := wipe.GetServerState(data, srv.UUID)

		// Check for error. Otherwise, break if we're not running.
		if err != nil {
			fmt.Println(err)
		} else {
			if state == "offline" {
				debug.SendDebugMsg(srv.UUID, data.DebugLevel, 4, "Found server offline. Continuing..")

				break
			}
		}

		// Increment i.
		i++

		// Kill the server after 15 seconds.
		if i == 15 {
			debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Found up for 15 seconds. Trying to kill server...")
			wipe.KillServer(data, srv.UUID)
		}

		// Give up after a minute.
		if i > 60 {
			debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Server halt timed out...")

			break
		}

		// Sleep every second to avoid unnecessary CPU cycles.
		time.Sleep(time.Duration(time.Second))
	}

	debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Processing files...")

	// Process and delete files.
	wipe.ProcessFiles(data, srv.UUID)

	debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Starting server back up...")

	// Start server back up.
	wipe.StartServer(data, srv.UUID)

	// Make sure the server starts back up.
	i = 0
	failed := 0

	for true {
		// Check if the server is running or starting. If confirmed, break loop.
		state, err := wipe.GetServerState(data, srv.UUID)

		// Check for error. Otherwise, break if we're not running.
		if err != nil {
			fmt.Println(err)
		} else {
			if state == "starting" || state == "running" {
				debug.SendDebugMsg(srv.UUID, data.DebugLevel, 4, "Found server starting/running. Continuing..")

				break
			}
		}

		// Increment i.
		i++

		// If we hit 15, start server again and reset i.
		if i == 15 {
			wipe.StartServer(data, srv.UUID)
			failed++

			i = 0
		}

		// Give up after 5 attempts.
		if failed > 5 {
			break
		}

		// Sleep every second to avoid unnecessary CPU cycles.
		time.Sleep(time.Duration(time.Second))
	}
}

func srv_handler(cfg *config.Config, srv *config.Server) error {
	var err error

	// Environmental overrides.
	wipe.EnvOverride(cfg, srv)

	// We need to retrieve the wipe data information first.
	var data wipe.Data

	// Process wipe data first.
	wipe.ProcessData(&data, cfg, srv)

	// Create cron job handler.
	c := cron.New()

	// If we have a single string, spawn a single cron job.
	for _, c_str := range data.CronStr {
		_, err = c.AddFunc("CRON_TZ="+data.TimeZone+" "+c_str, func() {
			wipe_server(cfg, srv, &data)

			next_wipe := "N/A"
			earliest := int64(1<<63 - 1)

			for _, cron := range c.Entries() {
				wipe_time := cron.Next.Unix()

				// If it's earlier, choose this one.
				//fmt.Println(strconv.FormatInt(wipe_time, 10) + " < " + strconv.FormatInt(earliest, 10))

				if wipe_time < earliest {
					earliest = wipe_time

					tz, err := time.LoadLocation(data.TimeZone)

					if err != nil {
						fmt.Println(err)

						continue
					}

					next_wipe = cron.Next.In(tz).Format("01-02-2006 3:04 PM MST")
				}
			}

			debug.SendDebugMsg(srv.UUID, data.DebugLevel, 1, "Server wiped. Next wipe date => "+next_wipe+".")
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	// Start cron job.
	c.Start()

	// Loop through each cron entry and print the next wipe date for debug.
	for _, job := range c.Entries() {
		tz, err := time.LoadLocation(data.TimeZone)

		if err != nil {
			fmt.Println(err)

			continue
		}

		// Retrieve the next time the job will be ran (Unix timestamp).
		next := job.Next.In(tz).Format("01-02-2006 3:04 PM MST")

		debug.SendDebugMsg(srv.UUID, data.DebugLevel, 1, "Next wipe date => "+next+".")
	}

	// See if we need to do a startup/first wipe.
	if srv.WipeFirst {
		wipe_server(cfg, srv, &data)
	}

	for true {
		// Loop through each cron entry.
		for _, job := range c.Entries() {
			// Retrieve the next time the job will be ran (Unix timestamp).
			now := time.Now().Unix()
			next := job.Next.Unix()

			// Loop through warning messages.
			for _, warning := range data.WarningMessages {
				wt := next - now

				// If what's remaining equals the warning time, we need to warn.
				if wt == int64(warning.WarningTime) {
					// Check if we're in running state.
					state, err := wipe.GetServerState(&data, srv.UUID)

					if err != nil {
						time.Sleep(time.Duration(time.Second))

						continue
					}

					if state != "running" {
						time.Sleep(time.Duration(time.Second))

						continue
					}

					warning_msg := *warning.Message
					format.FormatString(&warning_msg, int(wt))

					debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Sending warning message => "+warning_msg+".")

					err = wipe.SendMessage(&data, srv.UUID, warning_msg)

					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}

		// Sleep for one second to avoid unnecessary CPU cycles.
		time.Sleep(time.Duration(time.Second))
	}

	return err
}

func main() {
	var list bool
	var version bool
	var help bool

	// Setup simple flags (booleans).
	flag.BoolVar(&list, "list", false, "Print out config and exit.")
	flag.BoolVar(&list, "l", false, "Print out config and exit.")

	flag.BoolVar(&version, "version", false, "Print out version and exit.")
	flag.BoolVar(&version, "v", false, "Print out version and exit.")

	flag.BoolVar(&help, "help", false, "Print out help menu and exit.")
	flag.BoolVar(&help, "h", false, "Print out help menu and exit.")

	// Look for 'cfg' flag in command line arguments (default path: /etc/raw/raw.conf).
	configFile := flag.String("cfg", "/etc/raw/raw.conf", "The path to the Rust Auto Wipe config file.")

	// Parse flags.
	flag.Parse()

	// Check for version flag.
	if version {
		fmt.Print(VERSION)

		os.Exit(0)
	}

	// Check for help flag.
	if help {
		fmt.Print(HELP_MENU)

		os.Exit(0)
	}

	// Create config struct.
	cfg := config.Config{}

	// Set config defaults.
	cfg.SetDefaults()

	// Attempt to read config.
	err := cfg.LoadConfig(*configFile)

	// See if we want to automatically add servers.
	autoaddservers.AddServers(&cfg)

	// If we have no config, create the file with the defaults.
	if err != nil {
		// If there's an error and it contains "no such file", try to create the file with defaults.
		if strings.Contains(err.Error(), "no such file") {
			err = cfg.WriteDefaultsToFile(*configFile)

			if err != nil {
				fmt.Println("Failed to open config file and cannot create file.")
				fmt.Println(err)

				os.Exit(1)
			}
		}

		fmt.Println("WARNING - No config file found. Created config file at " + *configFile + " with defaults.")
	}

	// Check for list flag.
	if list {

		// Process environmental data before returning list.
		for i := 0; i < len(cfg.Servers); i++ {
			srv := &cfg.Servers[i]

			wipe.EnvOverride(&cfg, srv)
		}
		// Encode config as JSON string.
		json_data, err := json.MarshalIndent(cfg, "", "   ")

		if err != nil {
			fmt.Println(err)

			os.Exit(1)
		}

		fmt.Println(string(json_data))

		os.Exit(0)
	}

	// If we don't have any servers, what's the point?
	if len(cfg.Servers) < 1 {
		fmt.Println("[ERR] No servers found.")

		os.Exit(1)
	}

	// Loop through each server and execute Go routine.
	for _, srv := range cfg.Servers {
		// Check if we're enabled.
		if !srv.Enabled {
			continue
		}

		var srv_two config.Server = srv

		// Spawn Go routine.
		go srv_handler(&cfg, &srv_two)
	}

	// Signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc

	os.Exit(0)
}
