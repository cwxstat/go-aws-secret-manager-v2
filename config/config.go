package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func Config(region ...string) (aws.Config, error) {
	if len(region) == 0 {
		return config.LoadDefaultConfig(context.TODO())
	}
	return config.LoadDefaultConfig(context.TODO(), config.WithRegion(region[0]))
}
