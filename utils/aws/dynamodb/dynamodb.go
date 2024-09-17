package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	*dynamodb.Client
}

func NewDynamoDBClient() *DynamoDBClient {
	// Create a new DynamoDB client
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("failed to load configuration, " + err.Error())
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBClient{Client: client}
}

func (client *DynamoDBClient) AddItem(tableName string, item map[string]types.AttributeValue) error {
	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: &tableName,
	})

	return err
}
