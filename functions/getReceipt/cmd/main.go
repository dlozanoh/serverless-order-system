package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	sess      = session.Must(session.NewSession())
	db        = dynamodb.New(sess)
	s3Client  = s3.New(sess)
	tableName = os.Getenv("TABLE_NAME")
	bucket    = os.Getenv("BUCKET_NAME")
)

func generateSignedURL(key string, expiry time.Duration) (string, error) {
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(expiry)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	orderID := req.PathParameters["orderId"]
	res, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"orderId": {S: aws.String(orderID)},
		},
	})
	if err != nil || res.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Order not found",
		}, nil
	}

	key := fmt.Sprintf("order-receipts/%s.pdf", orderID)
	url, err := generateSignedURL(key, 5*time.Minute)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to generate signed URL",
		}, nil
	}

	body, _ := json.Marshal(map[string]string{
		"signedReceiptUrl": url,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
