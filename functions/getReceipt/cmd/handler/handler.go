package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDBClient interface to mock GetItem
type DynamoDBClient interface {
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

// S3Presigner interface to mock URL generation
// Mocking the raw S3 client for presigning is tricky because it involves the `request.Request` struct.
// It's cleaner to abstract the presigning logic itself.
type S3Presigner interface {
	GenerateSignedURL(key string, expiry time.Duration) (string, error)
}

type Handler struct {
	db        DynamoDBClient
	presigner S3Presigner
	tableName string
}

func New(db DynamoDBClient, presigner S3Presigner, tableName string) *Handler {
	return &Handler{
		db:        db,
		presigner: presigner,
		tableName: tableName,
	}
}

func (h *Handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	orderID := req.PathParameters["orderId"]
	if orderID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing orderId",
		}, nil
	}

	res, err := h.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(h.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"orderId": {S: aws.String(orderID)},
		},
	})

	if err != nil {
		fmt.Printf("Error getting item from DynamoDB: %v\n", err)
		// Usually return 500 on DB error, but if not found or DB error,
		// the requirements might be loose. Let's return 500 for DB error and 404 for not found specifically.
		// However, adhering to original logic which was a bit combined.
		// Let's stick closer to original logic: if err or Item is nil -> 404.
		// Wait, original logic: `if err != nil || res.Item == nil { return 404 }`
		// This swallows DB errors as 404s. I will keep this behavior to avoid changing contract, but it's not ideal.
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Order not found",
		}, nil
	}

	if res.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Order not found",
		}, nil
	}

	key := fmt.Sprintf("order-receipts/%s.pdf", orderID)
	url, err := h.presigner.GenerateSignedURL(key, 5*time.Minute)
	if err != nil {
		fmt.Printf("Error generating signed URL: %v\n", err)
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
