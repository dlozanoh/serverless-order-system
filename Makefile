.PHONY: all build deploy delete logs

APP_NAME=order-service
S3_BUCKET=sam-build-artifacts-bucket
STACK_NAME=order-service-stack
REGION=eu-west-3
EVENT_DIR=events

update:
	go mod tidy

clean:
	rm ./deployments

build:
	sam build -b deployments

deploy:
	sam deploy \
		--template-file deployments/template.yaml \
		--stack-name $(STACK_NAME) \
		--s3-prefix order-service-artifacts \
		--s3-bucket $(S3_BUCKET) \
		--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND \
		--region $(REGION) \
		--confirm-changeset \
		--parameter-overrides \
			AppName=$(APP_NAME)		

delete:
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME) \
		--region $(REGION) \

build-ReceiveOrderFunction:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -o bootstrap functions/receiveOrder/cmd/main.go
	cp bootstrap $(ARTIFACTS_DIR)/bootstrap

build-ProcessOrderFunction:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -o bootstrap functions/processOrder/cmd/main.go
	cp bootstrap $(ARTIFACTS_DIR)/bootstrap

build-GetReceiptFunction:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -o bootstrap functions/getReceipt/cmd/main.go
	cp bootstrap $(ARTIFACTS_DIR)/bootstrap

# add sam local invoke commands for testing
test-receive-order:
	sam local invoke --env-vars local.json receive-order --event $(EVENT_DIR)/receive_order_event.json -t deployments/receiveOrder/template.yaml

test-process-order:
	sam local invoke process-order --e $(EVENT_DIR)/process_order_event.json -t deployments/processOrder/template.yaml

test-get-receipt:
	sam local invoke get-receipt --e $(EVENT_DIR)/get_receipt_event.json -t deployments/getReceipt/template.yaml