package main

import (
	"log"
	"net/http"

	_ "github.com/phsaurav/echo_prod_blueprint/docs"
	"github.com/phsaurav/echo_prod_blueprint/internal/server"
)

// @title			JonoMot
// @version		0.1.0
// @description	A simple poll and voting API with user authentication
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Enter your JWT token directly (or optionally with 'Bearer ' prefix)
// @Security					BearerAuth
func main() {

	app, db, err := server.NewServer()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Start the HTTP server in a goroutine
	go func() {
		if err := app.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	server.GracefulShutdown(app, db, done)
	<-done
	log.Println("Server exiting")
}
