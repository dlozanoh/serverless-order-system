package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDynamoDB struct {
	mock.Mock
}

func (m *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

type MockPresigner struct {
	mock.Mock
}

func (m *MockPresigner) GenerateSignedURL(key string, expiry time.Duration) (string, error) {
	args := m.Called(key, expiry)
	return args.String(0), args.Error(1)
}

func TestHandle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB := new(MockDynamoDB)
		mockPresigner := new(MockPresigner)
		h := New(mockDB, mockPresigner, "test-table")

		mockDB.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.TableName == "test-table" && *input.Key["orderId"].S == "123"
		})).Return(&dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"orderId": {S: aws.String("123")},
			},
		}, nil).Once()

		mockPresigner.On("GenerateSignedURL", "order-receipts/123.pdf", 5*time.Minute).Return("https://signed.url", nil).Once()

		req := events.APIGatewayProxyRequest{
			PathParameters: map[string]string{"orderId": "123"},
		}

		resp, err := h.Handle(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.JSONEq(t, `{"signedReceiptUrl": "https://signed.url"}`, resp.Body)

		mockDB.AssertExpectations(t)
		mockPresigner.AssertExpectations(t)
	})

	t.Run("Order Not Found", func(t *testing.T) {
		mockDB := new(MockDynamoDB)
		mockPresigner := new(MockPresigner)
		h := New(mockDB, mockPresigner, "test-table")

		mockDB.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{
			Item: nil,
		}, nil).Once()

		req := events.APIGatewayProxyRequest{
			PathParameters: map[string]string{"orderId": "999"},
		}

		resp, err := h.Handle(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)

		mockDB.AssertExpectations(t)
		mockPresigner.AssertNotCalled(t, "GenerateSignedURL", mock.Anything, mock.Anything)
	})

	t.Run("DB Error", func(t *testing.T) {
		mockDB := new(MockDynamoDB)
		mockPresigner := new(MockPresigner)
		h := New(mockDB, mockPresigner, "test-table")

		mockDB.On("GetItem", mock.Anything).Return(nil, errors.New("db error")).Once()

		req := events.APIGatewayProxyRequest{
			PathParameters: map[string]string{"orderId": "123"},
		}

		resp, err := h.Handle(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode) // As per current logic

		mockDB.AssertExpectations(t)
	})

	t.Run("Presign Error", func(t *testing.T) {
		mockDB := new(MockDynamoDB)
		mockPresigner := new(MockPresigner)
		h := New(mockDB, mockPresigner, "test-table")

		mockDB.On("GetItem", mock.Anything).Return(&dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"orderId": {S: aws.String("123")},
			},
		}, nil).Once()

		mockPresigner.On("GenerateSignedURL", mock.Anything, mock.Anything).Return("", errors.New("presign failed")).Once()

		req := events.APIGatewayProxyRequest{
			PathParameters: map[string]string{"orderId": "123"},
		}

		resp, err := h.Handle(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)

		mockDB.AssertExpectations(t)
		mockPresigner.AssertExpectations(t)
	})
}
