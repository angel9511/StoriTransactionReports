package utils

import (
	"StoriTransactionReports/internal/config"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func DownloadFromS3(bucket, key string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AWSRegion),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := s3.New(sess)

	output, err := svc.GetObjectWithContext(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to download file %s from S3 %s bucket: %v", key, bucket, err)
	}
	defer output.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, output.Body); err != nil {
		return "", fmt.Errorf("failed to read file content: %v", err)
	}

	return buf.String(), nil
}
