package routes

import (
	"StoriTransactionReports/internal/handlers"
	"StoriTransactionReports/internal/services"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	reportService := services.NewReportServiceImpl()
	reportHandler := handlers.NewReportHandler(reportService)

	mux.HandleFunc("/processReport", reportHandler.ProcessReportHandler)
}
