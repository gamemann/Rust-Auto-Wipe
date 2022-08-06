package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/gamemann/Rust-Auto-Wipe/internal/autoaddservers"
	"github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/internal/wipe"
	"github.com/gamemann/Rust-Auto-Wipe/pkg/debug"
	cron "github.com/robfig/cron/v3"
)

func wipe_server(cfg *config.Config, srv *config.Server, data *wipe.Data) {
	debug.SendDebugMsg(srv.UUID, data.DebugLevel, 1, "Wiping server...")

	// Process map seeds.
	if data.ChangeMapSeeds {
		debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Processing seeds...")

		wipe.ProcessSeeds(data, srv.UUID)
	}

	// Process host name.
	if data.ChangeHostName {
		debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Processing hostname...")

		wipe.ProcessHostName(data, srv.UUID)
	}

	debug.SendDebugMsg(srv.UUID, data.DebugLevel, 2, "Stopping server...")

	// We should stop the server (To Do: Implement something to check if server is running and force kill if so).
	wipe.StopServer(data, srv.UUID)

	// Wait until the server is confirmed stopped.
	i := 0

	for true {
		// Check if the server is running and when it is confirmed stop, break the loop.
		running, err := wipe.IsServerRunning(data, srv.UUID)

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
}

func srv_handler(cfg *config.Config, srv *config.Server) error {
	var err error

	// We need to retrieve the wipe data information first.
	var data wipe.Data

	// Process wipe data first.
	wipe.ProcessData(&data, cfg, srv)

	// Create cron job handler.
	c := cron.New()

	// If we have a single string, spawn a single cron job.
	if reflect.TypeOf(data.CronStr).String() == "string" {
		_, err := c.AddFunc(reflect.ValueOf(data.CronStr).String(), func() {
			wipe_server(cfg, srv, &data)
		})

		if err != nil {
			return err
		}

	} else if reflect.TypeOf(data.CronStr).String() == "[]interface {}" {
		var tmp []reflect.Value
		for _, cron := range reflect.ValueOf(data.CronStr).CallSlice(tmp) {
			fmt.Println("Doing " + cron.String())
		}
	}

	// See if we need to do a startup/first wipe.
	if srv.WipeFirst {
		wipe_server(cfg, srv, &data)
	}

	return err
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

	// See if we want to automatically add servers.
	if cfg.AutoAddServers {
		autoaddservers.AddServers(&cfg)
	}

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

	// If we don't have any servers, what's the point?
	if len(cfg.Servers) < 1 {
		fmt.Println("[ERR] No servers found.")

		return
	}

	// Loop through each server and execute Go routine.
	for _, srv := range cfg.Servers {
		// Check if we're enabled.
		if !srv.Enabled {
			continue
		}

		// Spawn Go routine.
		go srv_handler(&cfg, &srv)
	}

	// Signal.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
}
