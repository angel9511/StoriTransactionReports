package db

import (
	"StoriTransactionReports/internal/utils"
	"fmt"
)

type MockDatabase struct {
	PersistedTransactions []utils.Transaction
	ShouldFail            bool
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		PersistedTransactions: []utils.Transaction{},
	}
}

func (m *MockDatabase) BatchPersistTransactions(transactions []utils.Transaction) error {
	if m.ShouldFail {
		return fmt.Errorf("expected database error")
	}
	m.PersistedTransactions = append(m.PersistedTransactions, transactions...)
	return nil
}
