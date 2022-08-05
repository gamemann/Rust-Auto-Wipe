package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/gamemann/Rust-Auto-Wipe/internal/config"
	"github.com/gamemann/Rust-Auto-Wipe/internal/processor"
)

func srv_handler(cfg *config.Config, srv *config.Server, idx int) {
	// We need to retrieve the wipe data information first.
	var wipedata processor.WipeData

	wipedata.ProcessData(cfg, srv)

	// Create a repeating loop until the two signals are called.
	for true {
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
