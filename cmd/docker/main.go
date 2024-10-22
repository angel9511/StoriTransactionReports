package main

import (
	"StoriTransactionReports/internal/config"
	"StoriTransactionReports/internal/routes"
	"log"
	"net/http"
)

func main() {
	config.Init()
	defer config.CloseDatabaseConnection()

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
