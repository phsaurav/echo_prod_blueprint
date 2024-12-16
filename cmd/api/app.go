package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	repository "github.com/phsaurav/go_echo_base/cmd"
	"github.com/phsaurav/go_echo_base/config"
	"github.com/phsaurav/go_echo_base/pkg/logger"
)

type application struct {
	repo   repository.Storage
	config config.Config
	log    *logger.Logger
}

func (app *application) run(mux *echo.Echo) error {
	// Docs
	//docs.SwaggerInfo.Version = version
	//docs.SwaggerInfo.Host = application.config.apiURL
	//docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// Start server
	go func() {
		if err := mux.StartServer(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.log.Fatal("Failed to start server")
		}
	}()

	app.log.Infof("Server started at %s in %s environment", app.config.Addr, app.config.Env)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.log.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mux.Shutdown(ctx); err != nil {
		app.log.Errorf("Server forced to shutdown: %v", err)
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	app.log.Info("Server exited gracefully")
	return nil
}
