package handlers

import (
	"StoriTransactionReports/internal/services"
	"StoriTransactionReports/internal/utils"
	"fmt"
	"io"
	"net/http"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) ProcessReportHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	transactions, err := utils.ParseCSV(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read CSV: %v", err), http.StatusBadRequest)
		return
	}

	recipient := r.Header.Get("recipient")
	if recipient == "" {
		http.Error(w, fmt.Sprintf("request recipient header missing"), http.StatusBadRequest)
		return
	}

	err = h.service.ProcessReport(transactions, recipient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to process report: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
