package main

import (
	"StoriTransactionReports/internal/config"
	"StoriTransactionReports/internal/db"
	"StoriTransactionReports/internal/handlers"
	"StoriTransactionReports/internal/services"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func main() {
	config.Init()

	postgresDB, err := db.NewPostgresDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer postgresDB.Close()
	log.Println("Starting Lambda handler...")
	reportService := services.NewReportServiceImpl(postgresDB)
	lambdaHandler := handlers.NewReportLambdaHandler(reportService)
	lambda.Start(lambdaHandler.HandleLambdaEvent)
}
