package libs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func FetchSecret(ctx context.Context, smClient *secretsmanager.Client, arn string) (string, error) {
	if arn == "" {
		// Use log package here since logger might not be initialized yet
		return "", fmt.Errorf("secret ARN cannot be empty")
	}
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(arn),
	}
	result, err := smClient.GetSecretValue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret %s: %w", arn, err)
	}
	if result.SecretString == nil {
		return "", fmt.Errorf("secret value for %s is nil or binary, expected string", arn)
	}
	return *result.SecretString, nil
}
