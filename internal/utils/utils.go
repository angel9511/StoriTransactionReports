package utils

import (
	"StoriTransactionReports/internal/config"
	"bytes"
	"encoding/csv"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"strconv"
	"strings"
)

var (
	SendEmail = SendEmailFunc
)

func ParseCSV(data []byte) ([]Transaction, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	var transactions []Transaction

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV file: %v", err)
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse ID: %v", err)
		}

		date := record[1]

		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Amount: %v", err)
		}

		transaction := Transaction{
			ID:     id,
			Date:   date,
			Amount: amount,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func SendEmailFunc(recipient, subject, body string) error {
	log.Print("Sending email to ", recipient)
	formattedBody := strings.ReplaceAll(body, "\n", "<br>")

	emailBody := fmt.Sprintf(`
		<html>
			<body>
				<div style="font-family: Arial, sans-serif; padding: 20px; background-color: #f9f0ff; border-radius: 15px;">
					<div style="border: 2px solid #e5d1fa; padding: 20px; border-radius: 10px; background-color: #ffffff; color: #6b5b95;">
						%s
					</div><br><br>
					<img src="cid:logo" alt="Logo" style="display: block; margin: 0 auto; width: 150px;">
				</div>
			</body>
		</html>
    `, formattedBody)

	m := gomail.NewMessage()
	m.SetHeader("From", config.SenderEmail)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	logoPath := "assets/logo.png"
	m.Embed(logoPath, gomail.SetHeader(map[string][]string{
		"Content-ID": {"<logo>"},
	}))

	d := gomail.NewDialer(config.SmtpHost, config.SmtpPort, config.SenderEmail, config.SenderPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Println("Email Sent Successfully.")
	return nil
}
