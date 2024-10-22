package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

var (
	SenderEmail    string
	SenderPassword string
	SmtpHost       string
	SmtpPort       int
	DbHost         string
	DbPort         int
	DbUser         string
	DbPassword     string
	DbName         string
	AWSRegion      string
)

var DB *sql.DB

func Init() {
	log.Println("Fetching environment variables...")
	initMailingConfig()
	initDbConfig()
	initAWSConfig()

	log.Println("Initializing database connection...")
	initializeDatabaseConnection()

	log.Println("Initialization complete")
}

func initMailingConfig() {
	SenderEmail = os.Getenv("SENDER_EMAIL")
	SenderPassword = os.Getenv("SENDER_PASSWORD")
	SmtpHost = os.Getenv("SMTP_HOST")
	SmtpPort = 587

	if SenderEmail == "" || SenderPassword == "" || SmtpHost == "" {
		panic("Mailing config is not properly set. Please check your environment variables.")
	}
}

func initDbConfig() {
	DbHost = os.Getenv("DB_HOST")
	DbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbName = os.Getenv("DB_NAME")

	if DbHost == "" || DbPort == 0 || DbUser == "" || DbName == "" || DbPassword == "" {
		panic("Database config is not properly set. Please check your environment variables.")
	}
}

func initAWSConfig() {
	AWSRegion = os.Getenv("AWS_REGION")

	if AWSRegion == "" {
		panic("AWS Config is not properly set. Please check your environment variables.")
	}
}

func initializeDatabaseConnection() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		DbHost,
		DbPort,
		DbUser,
		DbPassword,
		DbName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Database connection established successfully.")
}

func CloseDatabaseConnection() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed.")
	}
}
