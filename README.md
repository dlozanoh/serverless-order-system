# Serverless Order Processing System (Go + AWS SAM)

### Stack:
- AWS Lambda (Go)
- API Gateway
- SQS
- DynamoDB
- S3 (PDF Receipts)
- SAM for IaC

### Setup

1. Install [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
2. Build:
    make build
3. Deploy:
    make sam-deploy
4. Test:
- POST /orders
- GET /receipt/{orderId}


