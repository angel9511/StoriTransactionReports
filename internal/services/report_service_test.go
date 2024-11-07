package services

import (
	"StoriTransactionReports/internal/db"
	"StoriTransactionReports/internal/utils"
	"fmt"
	"testing"
)

func TestProcessReport(t *testing.T) {
	tests := []struct {
		name          string
		transactions  []utils.Transaction
		recipient     string
		dbSetup       func(mockDB *db.MockDatabase)
		expectedBody  string
		expectedError string
		mockSendEmail func(expectedBody, expectedRecipient string)
	}{
		{
			name: "Successful test case",
			transactions: []utils.Transaction{
				{ID: 35, Date: "7/15", Amount: 60.5},
				{ID: 36, Date: "7/28", Amount: -10.3},
				{ID: 37, Date: "8/2", Amount: -20.46},
				{ID: 39, Date: "8/13", Amount: 10.0},
			},
			recipient: "adavilam@unal.edu.co",
			dbSetup: func(mockDB *db.MockDatabase) {
				mockDB.ShouldFail = false
			},
			expectedError: "",
			expectedBody: `Total balance is 39.74

Number of transactions in July: 2
Number of transactions in August: 2
Average debit amount: -15.38
Average credit amount: 35.25`,
			mockSendEmail: func(expectedBody, expectedRecipient string) {
				utils.SendEmail = func(recipient, subject, body string) error {
					if recipient != expectedRecipient {
						return fmt.Errorf("unexpected recipient, got %s wanted %s", recipient, expectedRecipient)
					}
					if body != expectedBody {
						return fmt.Errorf("unexpected body, got %s wanted %s", body, expectedBody)
					}
					return nil
				}
			},
		},
		{
			name: "Failed to persist transactions",
			transactions: []utils.Transaction{
				{ID: 35, Date: "7/15", Amount: 60.5},
			},
			recipient: "adavilam@unal.edu.co",
			dbSetup: func(mockDB *db.MockDatabase) {
				mockDB.ShouldFail = true
			},
			expectedError: "failed to persist transactions: expected database error",
			expectedBody:  "",
			mockSendEmail: func(expectedBody, expectedRecipient string) {},
		},
		{
			name: "Failed to send email",
			transactions: []utils.Transaction{
				{ID: 35, Date: "7/15", Amount: 60.5},
				{ID: 36, Date: "7/28", Amount: -10.3},
				{ID: 37, Date: "8/2", Amount: -20.46},
				{ID: 39, Date: "8/13", Amount: 10.0},
			},
			recipient: "adavilam@unal.edu.co",
			dbSetup: func(mockDB *db.MockDatabase) {
				mockDB.ShouldFail = false
			},
			expectedError: "failed to send email: email error",
			expectedBody: `Total balance is 39.74

Number of transactions in July: 2
Number of transactions in August: 2
Average debit amount: -15.38
Average credit amount: 35.25`,
			mockSendEmail: func(expectedBody, expectedRecipient string) {
				utils.SendEmail = func(recipient, subject, body string) error {
					if recipient != expectedRecipient {
						return fmt.Errorf("unexpected recipient, got %s wanted %s", recipient, expectedRecipient)
					}
					if body != expectedBody {
						return fmt.Errorf("unexpected body, got %s wanted %s", body, expectedBody)
					}
					return fmt.Errorf("email error")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := db.NewMockDatabase()
			tc.dbSetup(mockDB)

			tc.mockSendEmail(tc.expectedBody, tc.recipient)
			svc := NewReportServiceImpl(mockDB)

			err := svc.ProcessReport(tc.transactions, tc.recipient)

			if tc.expectedError == "" && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tc.expectedError != "" && err.Error() != tc.expectedError {
				t.Errorf("expected error %s, got %v", tc.expectedError, err)
			}

			if tc.expectedError == "" && len(mockDB.PersistedTransactions) != len(tc.transactions) {
				t.Errorf("expected %d transactions to be persisted, but got %d", len(tc.transactions), len(mockDB.PersistedTransactions))
			}
		})
	}
}
