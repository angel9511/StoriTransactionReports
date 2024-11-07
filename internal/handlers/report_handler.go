package handlers

import (
	"StoriTransactionReports/internal/services"
	"StoriTransactionReports/internal/utils"
	"encoding/json"
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
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload utils.SummaryRequestPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	if payload.Recipient == "" {
		http.Error(w, "request missing recipient", http.StatusBadRequest)
		return
	}

	if payload.Transactions == "" {
		http.Error(w, "request missing transactions", http.StatusBadRequest)
		return
	}

	transactions, err := utils.ParseCSV([]byte(payload.Transactions))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse CSV: %v", err), http.StatusBadRequest)
		return
	}

	err = h.service.ProcessReport(transactions, payload.Recipient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to process report: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("transactions processed successfully"))
}
