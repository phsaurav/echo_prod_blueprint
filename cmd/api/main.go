package main

import (
	"log"
	"net/http"

	"github.com/phsaurav/go_echo_base/internal/server"
)

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
