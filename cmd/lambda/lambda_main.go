package main

import (
	"StoriTransactionReports/internal/config"
	"StoriTransactionReports/internal/handlers"
	"StoriTransactionReports/internal/services"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func main() {
	config.Init()
	defer config.CloseDatabaseConnection()

	log.Println("Starting Lambda handler...")
	reportService := services.NewReportServiceImpl()
	lambdaHandler := handlers.NewReportLambdaHandler(reportService)
	lambda.Start(lambdaHandler.HandleLambdaEvent)
}
