package sns

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPublishEvent_Success(t *testing.T) {
	mockPublishFunc := func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
		assert.Equal(t, "arn:aws:sns:us-east-1:123456789012:TestTopic", *params.TopicArn)
		assert.JSONEq(t, `{"event_type":"UserSignup","event_data":{"email":"test@example.com"},"timestamp":"2025-01-15T10:00:00Z"}`, *params.Message)
		return &sns.PublishOutput{
			MessageId: aws.String("12345"),
		}, nil
	}

	client := &AwsSnsClient{
		publishFunc: mockPublishFunc,
	}

	event := Event{
		EventType: "UserSignup",
		EventData: map[string]interface{}{
			"email": "test@example.com",
		},
		Timestamp: "2025-01-15T10:00:00Z",
	}

	// Act
	err := client.PublishEvent(context.Background(), "arn:aws:sns:us-east-1:123456789012:TestTopic", event)

	// Assert
	assert.NoError(t, err)
}

func TestPublishEvent_Error(t *testing.T) {
	// Arrange
	mockPublishFunc := func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
		return nil, errors.New("mock publish error")
	}

	client := &AwsSnsClient{
		publishFunc: mockPublishFunc,
	}

	event := Event{
		EventType: "UserSignup",
		EventData: map[string]interface{}{
			"email": "example@example.com",
		},
		Timestamp: "2025-01-15T10:00:00Z",
	}

	// Act
	err := client.PublishEvent(context.Background(), "arn:aws:sns:us-east-1:123456789012:YourTopicName", event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "failed to publish event, mock publish error", err.Error())
}
