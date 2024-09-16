package parameterstore

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type SSMClient struct {
	*ssm.Client
}

func NewSSMClient() *SSMClient {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("failed to load configuration, " + err.Error())
	}

	// Create an SSM client
	client := ssm.NewFromConfig(cfg)
	return &SSMClient{client}
}

func (client *SSMClient) GetParameterValue(paramName string, decrypt bool) (string, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(decrypt),
	}

	result, err := client.GetParameter(context.Background(), input)
	if err != nil {
		return "", err
	}

	return *result.Parameter.Value, nil
}

func (client *SSMClient) SetParameter(name, value string, overwrite bool) error {
	input := &ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Overwrite: aws.Bool(overwrite),
		Type:      types.ParameterTypeSecureString, // Set the type as needed (String, SecureString, etc.)
	}

	_, err := client.PutParameter(context.Background(), input)
	return err
}
