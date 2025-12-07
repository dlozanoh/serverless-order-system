package main

import (
	"os"
	"time"

	"serverless-order-system/functions/getReceipt/cmd/handler"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3PresignerImpl struct {
	client *s3.S3
	bucket string
}

func (s *S3PresignerImpl) GenerateSignedURL(key string, expiry time.Duration) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return req.Presign(expiry)
}

var (
	sess      = session.Must(session.NewSession())
	db        = dynamodb.New(sess)
	s3Client  = s3.New(sess)
	tableName = os.Getenv("TABLE_NAME")
	bucket    = os.Getenv("BUCKET_NAME")

	presigner = &S3PresignerImpl{client: s3Client, bucket: bucket}
	h         = handler.New(db, presigner, tableName)
)

func main() {
	lambda.Start(h.Handle)
}
