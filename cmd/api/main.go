package main

import (
	"github.com/phsaurav/go_echo_base/config"
	"github.com/phsaurav/go_echo_base/pkg/logger"
)

const version = "0.1.0"

func main() {
	// Load the application configuration from the specified directory.
	cfg, err := config.LoadConfig("config")
	if err != nil {
		// If an error occurs while loading the configuration, panic with the error.
		panic(err)
	}

	log := logger.NewLogger()
	log.SetLevel(cfg.LogLevel)

	mainApp := &app{
		config: cfg,
		log:    log,
	}

	mux := mainApp.mount()

	if err := mainApp.run(mux); err != nil {
		log.Fatal().Err(err).Msg("Application failed to run")
	}
}
