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
	defer func(log *logger.Logger) {
		err := log.Sync()
		if err != nil {

		}
	}(log)

	app := &appStruct{
		config: cfg,
		log:    log,
	}

	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Error("Application failed to run")
	}
}
