.PHONY: all build deploy delete logs

APP_NAME=order-service
S3_BUCKET=sam-build-artifacts-bucket
STACK_NAME=order-service-stack
REGION=eu-west-3
EVENT_DIR=testdata/events

update:
	go mod tidy

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
	GOOS=linux go build -o bootstrap ./cmd/receive_order/main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.

build-ProcessOrderFunction:
	GOOS=linux go build -o bootstrap ./cmd/process_order/main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.

build-GetReceiptFunction:
	GOOS=linux go build -o bootstrap ./cmd/get_receipt/main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.
	
# add sam local invoke commands for testing
test-receive-order:
	sam local invoke receive-order --event $(EVENT_DIR)/receive_order_event.json

test-process-order:
	sam local invoke process-order --event $(EVENT_DIR)/process_order_event.json

test-get-receipt:
	sam local invoke get-receipt --event $(EVENT_DIR)/get_receipt_event.json