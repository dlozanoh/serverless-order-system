package handler

import (
	"context"
	"encoding/json"
	"log"

	sqsClient "serverless-order-system/internal/sqs"

	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandler struct {
	sqsClient sqsClient.Client
}

type Order struct {
	OrderID      string `json:"orderId"`
	CustomerName string `json:"customerName"`
}

func New(sqsClient sqsClient.Client) *LambdaHandler {
	return &LambdaHandler{
		sqsClient: sqsClient,
	}
}

func (lh *LambdaHandler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var order Order
	err := json.Unmarshal([]byte(req.Body), &order)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input"}, nil
	}
	body, _ := json.Marshal(order)
	log.Printf("Received order: %s", body)
	err = lh.sqsClient.SendMessage(ctx, string(body))

	if err != nil {
		log.Printf("Failed to send message to SQS: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to send message"}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 202, Body: "Order received"}, nil
}
