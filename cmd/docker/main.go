package main

import (
	"StoriTransactionReports/internal/config"
	"StoriTransactionReports/internal/db"
	"StoriTransactionReports/internal/routes"
	"log"
	"net/http"
)

func main() {
	config.Init()

	mux := http.NewServeMux()

	postgresDB, err := db.NewPostgresDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer postgresDB.Close()
	routes.RegisterRoutes(mux, postgresDB)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
