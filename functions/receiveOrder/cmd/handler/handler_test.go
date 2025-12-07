package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSQSClient is a mock implementation of the sqsClient.Client interface
type MockSQSClient struct {
	mock.Mock
}

func (m *MockSQSClient) SendMessage(ctx context.Context, message string) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful Order",
			requestBody:    `{"orderId": "123", "customerName": "John Doe"}`,
			mockReturnErr:  nil,
			expectedStatus: 202,
			expectedBody:   "Order received",
		},
		{
			name:           "Invalid JSON",
			requestBody:    `Invalid JSON`,
			mockReturnErr:  nil,
			expectedStatus: 400,
			expectedBody:   "Invalid input",
		},
		{
			name:           "SQS Error",
			requestBody:    `{"orderId": "123", "customerName": "John Doe"}`,
			mockReturnErr:  errors.New("sqs error"),
			expectedStatus: 500,
			expectedBody:   "Failed to send message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSqs := new(MockSQSClient)
			h := New(mockSqs)

			if tt.expectedStatus != 400 {
				mockSqs.On("SendMessage", mock.Anything, mock.Anything).Return(tt.mockReturnErr)
			}

			req := events.APIGatewayProxyRequest{
				Body: tt.requestBody,
			}

			resp, err := h.Handle(context.Background(), req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, resp.Body)

			if tt.expectedStatus != 400 {
				mockSqs.AssertExpectations(t)
			}
		})
	}
}
