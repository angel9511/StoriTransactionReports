package services

import (
	"StoriTransactionReports/internal/db"
	"StoriTransactionReports/internal/utils"
	"fmt"
	"log"
	"time"
)

type ReportService interface {
	ProcessReport(transactions []utils.Transaction, recipient string) error
}

type ReportServiceImpl struct {
	db db.Database
}

func NewReportServiceImpl(database db.Database) *ReportServiceImpl {
	return &ReportServiceImpl{db: database}
}

func (s *ReportServiceImpl) ProcessReport(transactions []utils.Transaction, recipient string) error {
	log.Printf("Starting to process transactions, transaction count: %d", len(transactions))

	if err := s.db.BatchPersistTransactions(transactions); err != nil {
		return fmt.Errorf("failed to persist transactions: %v", err)
	}

	summary := s.computeSummary(transactions)

	err := utils.SendEmail(recipient, "Transaction Summary", summary)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (s *ReportServiceImpl) computeSummary(transactions []utils.Transaction) string {
	totalBalance := 0.0
	monthlyTransactionCount := make(map[string]int)
	totalCredit, totalDebit := 0.0, 0.0
	numCredits, numDebits := 0, 0

	for _, tx := range transactions {
		parsedDate, _ := time.Parse("1/2", tx.Date)
		month := parsedDate.Month().String()
		monthlyTransactionCount[month]++

		totalBalance += tx.Amount
		if tx.Amount > 0 {
			totalCredit += tx.Amount
			numCredits++
		} else {
			totalDebit += tx.Amount
			numDebits++
		}
	}

	avgCredit := 0.0
	if numCredits > 0 {
		avgCredit = totalCredit / float64(numCredits)
	}

	avgDebit := 0.0
	if numDebits > 0 {
		avgDebit = totalDebit / float64(numDebits)
	}

	summary := fmt.Sprintf("Total balance is %.2f\n\n", totalBalance)
	for _, month := range utils.MONTHS {
		if count, exists := monthlyTransactionCount[month]; exists && count > 0 {
			summary += fmt.Sprintf("Number of transactions in %s: %d\n", month, count)
		}
	}
	summary += fmt.Sprintf("Average debit amount: %.2f\nAverage credit amount: %.2f", avgDebit, avgCredit)

	log.Printf("%s", summary)
	return summary
}
