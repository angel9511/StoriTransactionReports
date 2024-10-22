package db

import (
	"StoriTransactionReports/internal/config"
	"StoriTransactionReports/internal/utils"
	"fmt"
	"log"
	"strings"
)

func BatchPersistTransactions(transactions []utils.Transaction) error {
	log.Println("Starting batch persist transaction.")
	tx, err := config.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Printf("transaction rolled back due to error: %v", err)
		}
	}()

	sqlStr := "INSERT INTO transactions (id, date, amount) VALUES "
	vals := []interface{}{}
	const rowSQL = "($%d, $%d, $%d)"
	inserts := []string{}

	for i, transaction := range transactions {
		inserts = append(inserts, fmt.Sprintf(rowSQL, i*3+1, i*3+2, i*3+3))
		vals = append(vals, transaction.ID, transaction.Date, transaction.Amount)
	}

	sqlStr = sqlStr + strings.Join(inserts, ",") + " ON CONFLICT (id) DO NOTHING"

	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(vals...)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Println("Batch persist transaction complete.")
	return nil
}
