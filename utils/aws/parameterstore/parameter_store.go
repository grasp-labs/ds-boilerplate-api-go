package parameterstore

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func NewSSMClient() *ssm.Client {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("failed to load configuration, " + err.Error())
	}

	// Create an SSM client
	client := ssm.NewFromConfig(cfg)
	return client
}
