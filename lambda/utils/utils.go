package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	cachedJWTSecret string
	secretOnce      sync.Once
	cachedErr       error
)

// FetchJWTSecret loads and caches the JWT secret key (from key "key") in the secret JSON.
func FetchJWTSecret(secretArn string) (string, error) {
	secretOnce.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			cachedErr = fmt.Errorf("failed to load AWS config: %w", err)
			return
		}

		client := secretsmanager.NewFromConfig(cfg)
		resp, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
			SecretId:     &secretArn,
			VersionStage: aws.String("AWSCURRENT"),
		})
		if err != nil {
			cachedErr = fmt.Errorf("failed to get secret value: %w", err)
			return
		}

		var secretMap map[string]string
		err = json.Unmarshal([]byte(*resp.SecretString), &secretMap)
		if err != nil {
			cachedErr = fmt.Errorf("failed to unmarshal secret JSON: %w", err)
			return
		}

		secret, ok := secretMap["key"]
		if !ok {
			cachedErr = fmt.Errorf(`secret JSON does not contain "key"`)
			return
		}

		cachedJWTSecret = secret
	})

	return cachedJWTSecret, cachedErr
}
