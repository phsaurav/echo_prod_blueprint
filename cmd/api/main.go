package main

import (
	"database/sql"

	repository "github.com/phsaurav/go_echo_base/cmd"
	"github.com/phsaurav/go_echo_base/config"
	"github.com/phsaurav/go_echo_base/db"
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

	//Logger
	log := logger.NewLogger()
	log.SetLevel(cfg.LogLevel)
	defer func(log *logger.Logger) {
		err := log.Sync()
		if err != nil {
			log.Fatalf("Error syncing logger: %v", err)
		}
	}(log)

	//Main Database
	database, err := db.New(
		cfg.Db.Addr,
		cfg.Db.MaxOpenConns,
		cfg.Db.MaxIdleConns,
		cfg.Db.MaxIdleTime,
	)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(database)

	repo := repository.NewStorage(database)

	app := &application{
		repo:   repo,
		config: cfg,
		log:    log,
	}

	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Error("Application failed to run")
	}
}
