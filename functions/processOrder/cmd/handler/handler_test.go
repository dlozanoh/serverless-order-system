package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDynamoDB struct {
	mock.Mock
}

func (m *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(input)
	// Safely cast the first return value, handling potential nil or type mismatches if needed,
	// but here we expect *dynamodb.PutItemOutput.
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

type MockS3 struct {
	mock.Mock
}

func (m *MockS3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func TestHandle(t *testing.T) {
	t.Run("Successful Processing", func(t *testing.T) {
		mockDB := new(MockDynamoDB)
		mockS3 := new(MockS3)
		h := New(mockDB, mockS3, "test-table", "test-bucket")

		// Expect PutItem to be called
		mockDB.On("PutItem", mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil).Once()

		// Expect PutObject to be called
		mockS3.On("PutObject", mock.AnythingOfType("*s3.PutObjectInput")).Return(&s3.PutObjectOutput{}, nil).Once()

		sqsEvent := events.SQSEvent{
			Records: []events.SQSMessage{
				{
					Body: `{"orderId": "123", "customerName": "Test User"}`,
				},
			},
		}

		err := h.Handle(context.Background(), sqsEvent)
		assert.NoError(t, err)

		mockDB.AssertExpectations(t)
		mockS3.AssertExpectations(t)
	})

	t.Run("DynamoDB Error", func(t *testing.T) {
		mockDB := new(MockDynamoDB)
		mockS3 := new(MockS3)
		h := New(mockDB, mockS3, "test-table", "test-bucket")

		mockDB.On("PutItem", mock.Anything).Return(nil, errors.New("dynamo error")).Once()

		// S3 should NOT be called if DynamoDB fails (based on current implementation logic,
		// actually the current implementation continues on error? let's check code.
		// Ah, it has 'continue' so it skips the rest of the loop for that message. Correct.)

		sqsEvent := events.SQSEvent{
			Records: []events.SQSMessage{
				{
					Body: `{"orderId": "456", "customerName": "Test User 2"}`,
				},
			},
		}

		err := h.Handle(context.Background(), sqsEvent)
		assert.NoError(t, err) // Handler returns nil even on error, just logs

		mockDB.AssertExpectations(t)
		mockS3.AssertNotCalled(t, "PutObject", mock.Anything)
	})
}
