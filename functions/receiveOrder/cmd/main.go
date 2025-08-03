package main

import (
	"serverless-order-system/functions/receiveOrder/cmd/handler"
	"serverless-order-system/functions/receiveOrder/cmd/internal"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	provider      = internal.NewProvider()
	sqsClient     = provider.ProvideSqsClient()
	lambdaHandler = handler.New(sqsClient)
)

func main() {
	lambda.Start(lambdaHandler.Handle)
}
