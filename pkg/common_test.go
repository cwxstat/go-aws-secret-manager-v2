package pkg

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"testing"
	"time"
)

type mockCreateSecret func(ctx context.Context,
	params *secretsmanager.CreateSecretInput,
	optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)

func (m mockCreateSecret) CreateSecret(ctx context.Context,
	params *secretsmanager.CreateSecretInput,
	optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
	return m(ctx, params, optFns...)
}

func TestCreateSecret(t *testing.T) {
	cases := []struct {
		client       func(t *testing.T) SecretsManagerCreateSecretAPI
		name         string
		description  string
		secretString string
		expect       []byte
	}{
		{
			client: func(t *testing.T) SecretsManagerCreateSecretAPI {
				return mockCreateSecret(func(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
					t.Helper()
					if params.Name == nil {
						t.Errorf("expected name to be set")
					}
					if e, a := "example", *params.Name; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return &secretsmanager.CreateSecretOutput{
						ARN: aws.String("arn:aws:secretsmanager:us-west-2:123456789012:secret:example-123456"),
					}, nil
				})
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			_, err := CreateSecret(ctx, c.client(t), "example", "value", "description")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

// DELETE

type mockDeleteSecret func(ctx context.Context,
	params *secretsmanager.DeleteSecretInput,
	optFns ...func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error)

func (m mockDeleteSecret) DeleteSecret(ctx context.Context,
	params *secretsmanager.DeleteSecretInput,
	optFns ...func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error) {
	return m(ctx, params, optFns...)
}

func TestDeleteSecret(t *testing.T) {
	cases := []struct {
		client   func(t *testing.T) SecretsManagerDeleteSecretAPI
		secretId string
		expect   []byte
	}{
		{
			client: func(t *testing.T) SecretsManagerDeleteSecretAPI {
				return mockDeleteSecret(func(ctx context.Context, params *secretsmanager.DeleteSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error) {
					t.Helper()
					if params.SecretId == nil {
						t.Errorf("expected name to be set")
					}
					if e, a := "example", *params.SecretId; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return &secretsmanager.DeleteSecretOutput{
						ARN: aws.String("arn:aws:secretsmanager:us-west-2:123456789012:secret:example-123456"),
					}, nil
				})
			},
		},
	}

	for _, c := range cases {
		t.Run(c.secretId, func(t *testing.T) {
			ctx := context.TODO()
			err := DeleteSecret(ctx, c.client(t), "example")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}

}

// GET

type mockGetSecret func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)

func (m mockGetSecret) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return m(ctx, params, optFns...)
}

func TestGetSecret(t *testing.T) {
	cases := []struct {
		client   func(t *testing.T) SecretsManagerGetSecretAPI
		name     string
		secretId string
		expect   []byte
	}{
		{
			client: func(t *testing.T) SecretsManagerGetSecretAPI {
				return mockGetSecret(func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
					t.Helper()
					if params.SecretId == nil {
						t.Errorf("expected name to be set")
					}
					if e, a := "example", *params.SecretId; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return &secretsmanager.GetSecretValueOutput{
						ARN:          aws.String("arn:aws:secretsmanager:us-west-2:123456789012:secret:example-123456"),
						Name:         aws.String("example"),
						SecretString: aws.String("my secret ....pss..."),
					}, nil
				})
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			secret, err := GetSecret(ctx, c.client(t), "example")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if e, a := "my secret ....pss...", secret; e != a {
				t.Errorf("expected %v, got %v", e, a)
			}
		})
	}

}

// LIST
type mockListSecret func(ctx context.Context,
	params *secretsmanager.ListSecretsInput,
	optFns ...func(*secretsmanager.Options)) (*secretsmanager.ListSecretsOutput, error)

func (m mockListSecret) ListSecrets(ctx context.Context,
	params *secretsmanager.ListSecretsInput,
	optFns ...func(*secretsmanager.Options)) (*secretsmanager.ListSecretsOutput, error) {
	return m(ctx, params, optFns...)
}

func TestListSecrets(t *testing.T) {
	list := []types.SecretListEntry{}
	for i, v := range []string{"a", "b", "c"} {
		list = append(list, types.SecretListEntry{
			ARN:                    aws.String("arn:aws:secretsmanager:us-east-1:123456789012:secret:MyTestDatabaseSecret-ABC123"),
			CreatedDate:            aws.Time(time.Now()),
			DeletedDate:            nil,
			Description:            aws.String("My test database secret: " + v),
			KmsKeyId:               aws.String(fmt.Sprintf("key-id-%d", i)),
			LastAccessedDate:       aws.Time(time.Now()),
			LastChangedDate:        nil,
			LastRotatedDate:        nil,
			Name:                   aws.String("MyTestDatabaseSecret" + v),
			OwningService:          nil,
			PrimaryRegion:          nil,
			RotationLambdaARN:      nil,
			RotationRules:          nil,
			SecretVersionsToStages: nil,
			Tags:                   nil,
		})
	}
	cases := []struct {
		client func(t *testing.T) SecretsManagerListSecretAPI
		name   string
		expect []byte
	}{
		{
			client: func(t *testing.T) SecretsManagerListSecretAPI {
				return mockListSecret(func(ctx context.Context, params *secretsmanager.ListSecretsInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.ListSecretsOutput, error) {
					t.Helper()

					return &secretsmanager.ListSecretsOutput{
						SecretList: list,
					}, nil
				})
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			list, err := ListSecrets(ctx, c.client(t))
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if len(list) != 3 {
				t.Errorf("expected 3 secrets, got %d", len(list))
			}
		})
	}

}

// UPDATE
type mockUpdateSecret func(ctx context.Context,
	params *secretsmanager.UpdateSecretInput,
	optFns ...func(*secretsmanager.Options)) (*secretsmanager.UpdateSecretOutput, error)

func (m mockUpdateSecret) UpdateSecret(ctx context.Context,
	params *secretsmanager.UpdateSecretInput,
	optFns ...func(*secretsmanager.Options)) (*secretsmanager.UpdateSecretOutput, error) {
	return m(ctx, params, optFns...)
}

func TestUpdateSecret(t *testing.T) {
	cases := []struct {
		client       func(t *testing.T) SecretsManagerUpdateSecretAPI
		name         string
		secretId     string
		secretString string
		expect       []byte
	}{
		{
			client: func(t *testing.T) SecretsManagerUpdateSecretAPI {
				return mockUpdateSecret(func(ctx context.Context, params *secretsmanager.UpdateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.UpdateSecretOutput, error) {
					t.Helper()
					if params.SecretId == nil || params.SecretString == nil {
						t.Errorf("expected secretID and SecretString to be set")
					}
					if e, a := "example", *params.SecretId; e != a {
						t.Errorf("expected %v, got %v", e, a)
					}
					return &secretsmanager.UpdateSecretOutput{
						ARN: aws.String("arn:aws:secretsmanager:us-west-2:123456789012:secret:example-123456"),
					}, nil
				})
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			err := UpdateSecret(ctx, c.client(t), "example", "value")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}

}
