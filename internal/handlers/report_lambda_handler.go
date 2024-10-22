package handlers

import (
	"StoriTransactionReports/internal/services"
	"StoriTransactionReports/internal/utils"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/url"
	"strings"
)

type ReportLambdaHandler struct {
	service services.ReportService
}

func NewReportLambdaHandler(service services.ReportService) *ReportLambdaHandler {
	return &ReportLambdaHandler{service: service}
}

func (h *ReportLambdaHandler) HandleLambdaEvent(ctx context.Context, s3Event events.S3Event) error {
	bucket := s3Event.Records[0].S3.Bucket.Name
	key, err := url.QueryUnescape(s3Event.Records[0].S3.Object.Key)
	if err != nil {
		return fmt.Errorf("failed to decode S3 object key: %v", err)
	}

	csvData, err := utils.DownloadFromS3(bucket, key)
	if err != nil {
		return fmt.Errorf("failed to obtain CSV file from S3: %v", err)
	}

	transactions, err := utils.ParseCSV([]byte(csvData))
	if err != nil {
		return fmt.Errorf("failed to read CSV: %v", err)
	}
	recipient := strings.TrimSuffix(key, ".csv")

	err = h.service.ProcessReport(transactions, recipient)
	if err != nil {
		return fmt.Errorf("failed to process report: %v", err)
	}

	return err
}
