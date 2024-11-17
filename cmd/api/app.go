package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/phsaurav/go_echo_base/config"
	"github.com/phsaurav/go_echo_base/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	config config.Config
	log    *logger.Logger
}

func (a *app) run(mux *echo.Echo) error {
	// Docs
	//docs.SwaggerInfo.Version = version
	//docs.SwaggerInfo.Host = app.config.apiURL
	//docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         a.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// Start server
	go func() {
		if err := mux.StartServer(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatal("Failed to start server")
		}
	}()

	a.log.Infof("Server started at %s in %s environment", a.config.Addr, a.config.Env)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.log.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mux.Shutdown(ctx); err != nil {
		a.log.Errorf("Server forced to shutdown: %v", err)
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	a.log.Info("Server exited gracefully")
	return nil
}
