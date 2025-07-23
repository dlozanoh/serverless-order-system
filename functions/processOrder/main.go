package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jung-kurt/gofpdf"
)

type Order struct {
	OrderID      string `json:"orderId"`
	CustomerName string `json:"customerName"`
}

var (
	sess      = session.Must(session.NewSession())
	db        = dynamodb.New(sess)
	s3Client  = s3.New(sess)
	tableName = os.Getenv("TABLE_NAME")
	bucket    = os.Getenv("BUCKET_NAME")
)

func generatePDF(order Order) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Order Receipt")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Order ID: %s", order.OrderID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", order.CustomerName))
	pdf.Ln(8)
	pdf.Cell(40, 10, "Thank you for your order!")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, msg := range sqsEvent.Records {
		var order Order
		if err := json.Unmarshal([]byte(msg.Body), &order); err != nil {
			fmt.Println("Error unmarshalling order:", err)
			continue
		}

		// Store in DynamoDB
		av, err := dynamodbattribute.MarshalMap(order)
		if err != nil {
			fmt.Println("Error marshalling order:", err)
			continue
		}
		_, err = db.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      av,
		})
		if err != nil {
			fmt.Println("Error saving to DynamoDB:", err)
			continue
		}

		// Generate and upload PDF receipt
		pdfBytes, err := generatePDF(order)
		if err != nil {
			fmt.Println("Error generating PDF:", err)
			continue
		}
		key := fmt.Sprintf("order-receipts/%s.pdf", order.OrderID)
		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   bytes.NewReader(pdfBytes),
		})
		if err != nil {
			fmt.Println("Error uploading PDF to S3:", err)
			continue
		}

		fmt.Printf("Processed order %s\n", order.OrderID)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
