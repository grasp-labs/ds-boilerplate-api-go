package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Client struct {
	client   *sqs.Client
	queueURL string
}

type UsageMessageAttributes struct {
	TenantID       types.MessageAttributeValue
	ProductID      types.MessageAttributeValue
	MemoryMB       types.MessageAttributeValue
	StartTimestamp types.MessageAttributeValue
	EndTimestamp   types.MessageAttributeValue
	Workflow       types.MessageAttributeValue
}

func NewUsageMessageAttributes(tenantID, productID, memoryMB, startTimestamp, endTimestamp, workflow string) UsageMessageAttributes {
	return UsageMessageAttributes{
		TenantID: types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(tenantID),
		},
		ProductID: types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(productID),
		},
		MemoryMB: types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(memoryMB),
		},
		StartTimestamp: types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(startTimestamp),
		},
		EndTimestamp: types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(endTimestamp),
		},
		Workflow: types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(workflow),
		},
	}
}

func NewSQSClient(queueUrl string) *Client {
	// Create an SQS client
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("failed to load configuration, " + err.Error())
	}

	// Create an SQS client
	client := sqs.NewFromConfig(cfg)
	return &Client{
		client:   client,
		queueURL: queueUrl,
	}
}

func (s *Client) SendUsageMessage(ctx context.Context, messageBody string, messageAttributes UsageMessageAttributes) error {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    aws.String(s.queueURL),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"product_id":      messageAttributes.ProductID,
			"tenant_id":       messageAttributes.TenantID,
			"memory_mb":       messageAttributes.MemoryMB,
			"start_timestamp": messageAttributes.StartTimestamp,
			"end_timestamp":   messageAttributes.EndTimestamp,
			"workflow":        messageAttributes.Workflow,
		},
	}

	// Send message to SQS
	_, err := s.client.SendMessage(ctx, input)
	return err
}
