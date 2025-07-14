.PHONY: build

build:
	GOOS=linux GOARCH=amd64 go build -o cmd/receive_order/main ./cmd/receive_order
	GOOS=linux GOARCH=amd64 go build -o cmd/process_order/main ./cmd/process_order
	GOOS=linux GOARCH=amd64 go build -o cmd/get_receipt/main ./cmd/get_receipt

sam-build:
	sam build -b deployments

sam-deploy:
	sam deploy --no-confirm-changeset --stack-name order-service --capabilities CAPABILITY_IAM