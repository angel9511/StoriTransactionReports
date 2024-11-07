package db

import (
	"StoriTransactionReports/internal/config"
	"StoriTransactionReports/internal/utils"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Database interface {
	BatchPersistTransactions(transactions []utils.Transaction) error
}

type PostgresDatabase struct {
	DB *sql.DB
}

func NewPostgresDatabase() (*PostgresDatabase, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPassword,
		config.DbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	fmt.Println("Database connection established successfully.")
	return &PostgresDatabase{DB: db}, nil
}

func (p *PostgresDatabase) Close() {
	if p.DB != nil {
		p.DB.Close()
		fmt.Println("Database connection closed.")
	}
}

func (p *PostgresDatabase) BatchPersistTransactions(transactions []utils.Transaction) error {
	log.Println("Starting batch persist transaction.")
	tx, err := p.DB.Begin()
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
