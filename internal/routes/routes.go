package routes

import (
	"StoriTransactionReports/internal/db"
	"StoriTransactionReports/internal/handlers"
	"StoriTransactionReports/internal/services"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, db db.Database) {
	reportService := services.NewReportServiceImpl(db)
	reportHandler := handlers.NewReportHandler(reportService)

	mux.HandleFunc("/processReport", reportHandler.ProcessReportHandler)
}
