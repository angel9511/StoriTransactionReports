package services

import (
	"StoriTransactionReports/internal/config"
	"StoriTransactionReports/internal/utils"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestProcessReport(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	config.DB = db

	tests := []struct {
		name          string
		transactions  []utils.Transaction
		recipient     string
		dbSetup       func()
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
			dbSetup: func() {
				// Mocking the batch insert as a single transaction
				mock.ExpectBegin() // Start the transaction

				// The expected SQL batch insert
				mock.ExpectPrepare("INSERT INTO transactions \\(id, date, amount\\) VALUES (.+) ON CONFLICT \\(id\\) DO NOTHING").
					ExpectExec().
					WithArgs(
						35, "7/15", 60.5,
						36, "7/28", -10.3,
						37, "8/2", -20.46,
						39, "8/13", 10.0,
					).
					WillReturnResult(sqlmock.NewResult(4, 4)) // Mock successful insert for 4 rows

				mock.ExpectCommit() // Expect a commit
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
			dbSetup: func() {
				mock.ExpectBegin() // Start the transaction

				mock.ExpectPrepare("INSERT INTO transactions \\(id, date, amount\\) VALUES (.+) ON CONFLICT \\(id\\) DO NOTHING").
					ExpectExec().
					WithArgs(35, "7/15", 60.5).
					WillReturnError(fmt.Errorf("db insert error")) // Simulate an insert error

				mock.ExpectRollback() // Expect a rollback due to error
			},
			expectedError: "failed to persist transactions: failed to execute batch insert: db insert error",
			expectedBody:  "",
			mockSendEmail: func(expectedBody, expectedRecipient string) {
			},
		},
		{
			name: "Failed to send email",
			transactions: []utils.Transaction{
				{ID: 35, Date: "7/15", Amount: 60.5},
			},
			recipient: "adavilam@unal.edu.co",
			dbSetup: func() {
				mock.ExpectBegin() // Start the transaction

				mock.ExpectPrepare("INSERT INTO transactions \\(id, date, amount\\) VALUES (.+) ON CONFLICT \\(id\\) DO NOTHING").
					ExpectExec().
					WithArgs(35, "7/15", 60.5).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Successful insert

				mock.ExpectCommit() // Commit the transaction
			},
			expectedError: "failed to send email: email error",
			expectedBody: `Total balance is 60.50

Number of transactions in July: 1
Average debit amount: 0.00
Average credit amount: 60.50`,
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
			tc.mockSendEmail(tc.expectedBody, tc.recipient)
			tc.dbSetup()
			svc := NewReportServiceImpl()

			err := svc.ProcessReport(tc.transactions, tc.recipient)

			if tc.expectedError == "" && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tc.expectedError != "" && err.Error() != tc.expectedError {
				t.Errorf("expected error %s, got %v", tc.expectedError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unmet database expectations: %v", err)
			}
		})
	}
}
