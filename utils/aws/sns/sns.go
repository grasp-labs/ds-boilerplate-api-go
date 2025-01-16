package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type AwsSnsClient struct {
	client      *sns.Client
	publishFunc func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type Event struct {
	EventType string                 `json:"event_type"`
	EventData map[string]interface{} `json:"event_data"`
	Timestamp string                 `json:"timestamp"`
}

func NewAwsSnsClient() (*AwsSnsClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %v", err)
	}
	client := sns.NewFromConfig(cfg)
	return &AwsSnsClient{
		client:      client,
		publishFunc: client.Publish,
	}, nil
}

func (c *AwsSnsClient) PublishEvent(ctx context.Context, topicArn string, event Event) error {
	if topicArn == "" {
		return fmt.Errorf("topicArn is required")
	}

	// Serialize the event to JSON
	eventStr, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event to JSON, %v", err)
	}

	_, err = c.publishFunc(ctx, &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(string(eventStr)),
	})

	if err != nil {
		return fmt.Errorf("failed to publish event, %v", err)
	}
	return nil
}
