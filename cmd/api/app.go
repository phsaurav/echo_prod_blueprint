package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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

var configCORS = echoMiddleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPatch,
	},
}

func (a *app) mount() *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(configCORS))

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to the Echo API!")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World 👋")
	})

	return e
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
		if err := mux.StartServer(srv); err != nil && err != http.ErrServerClosed {
			a.log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	a.log.Info().Str("addr", a.config.Addr).Str("env", a.config.Env).Msg("Server started")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.log.Info().Msg("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mux.Shutdown(ctx); err != nil {
		a.log.Error().Err(err).Msg("Server forced to shutdown")
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	a.log.Info().Msg("Server exited gracefully")
	return nil
}
