package internal

import (
	"context"

	sqsClient "serverless-order-system/internal/sqs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/kelseyhightower/envconfig"
)

type Provider struct {
	awsConfig    aws.Config
	envVariables *EnvVariables
}

func (p *Provider) ProvideSqsClient() sqsClient.Client {
	awsSqsClient := sqs.NewFromConfig(p.awsConfig)
	return sqsClient.New(awsSqsClient, p.envVariables.QueueURL)
}

type EnvVariables struct {
	Application string `envconfig:"APPLICATION" required:"true"`
	QueueURL    string `envconfig:"QUEUE_URL" required:"true"`
}

func NewProvider() *Provider {
	var envVars EnvVariables
	if err := envconfig.Process("", &envVars); err != nil {
		panic("Failed to load environment variables: " + err.Error())
	}
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("Failed to load AWS config: " + err.Error())
	}
	return &Provider{
		awsConfig:    awsConfig,
		envVariables: &envVars,
	}
}
