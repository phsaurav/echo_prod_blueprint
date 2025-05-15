package server

import (
	"context"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"

	"github.com/phsaurav/echo_prod_blueprint/config"
	"github.com/phsaurav/echo_prod_blueprint/internal/database"
	"github.com/phsaurav/echo_prod_blueprint/pkg/logger"
)

type Server struct {
	store  Store
	config config.Config
	log    *logger.Logger
	e      *echo.Echo
}

func NewServer() (*http.Server, database.Service, error) {
	// Load the application configuration from the specified directory.
	cfg, err := config.LoadConfig()
	if err != nil {
		// If an error occurs while loading the configuration, panic with the error.
		panic(err)
	}

	//Logger
	log := logger.NewLogger()
	log.SetLevel(cfg.LogLevel)

	db, err := database.New(
		cfg.Db.Addr,
		cfg.Db.MaxOpenConns,
		cfg.Db.MaxIdleConns,
		cfg.Db.MaxIdleTime,
	)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil, nil, err
	}

	store := NewStore(db)

	NewServer := &Server{
		store:  store,
		config: cfg,
		log:    log,
	}

	// Declare Server config
	app := &http.Server{
		Addr:         NewServer.config.Addr,
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  NewServer.config.IdleTimeout,
		ReadTimeout:  NewServer.config.ReadTimeout,
		WriteTimeout: NewServer.config.WriteTimeout,
	}

	return app, db, nil
}

func GracefulShutdown(apiServer *http.Server, db database.Service, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.

	var logging = logger.NewLogger()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	logging.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		logging.Errorf("Server forced to shutdown with error: %v", err)
	}

	// Close the database connection gracefully
	if err := db.Close(); err != nil {
		logging.Info("Error closing the database connection: %v", err)
	}

	// Flush any buffered logs.
	if err := logging.Sync(); err != nil && !strings.Contains(err.Error(), "invalid argument") {
		logging.Errorf("Error syncing logger: %v", err)
	}

	logging.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
