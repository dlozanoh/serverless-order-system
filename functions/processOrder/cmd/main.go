package main

import (
	"os"

	"serverless-order-system/functions/processOrder/cmd/handler"

	"github.com/aws/aws-lambda-go/lambda"
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
	h         = handler.New(db, s3Client, tableName, bucket)
)

func main() {
	lambda.Start(h.Handle)
}
