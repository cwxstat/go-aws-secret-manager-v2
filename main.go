package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/google/uuid"

	"github.com/cwxstat/go-aws-secret-manager-v2/config"
	"github.com/cwxstat/go-aws-secret-manager-v2/pkg"
)

func main() {

	var secretArn string
	var secretName string
	var value string

	cfg, err := config.Config("us-east-1")

	if err != nil {
		panic("Couldn't load config!")
	}

	secretName = "sampleSecret" + uuid.NewString()

	value = `some secret value`

	conn := secretsmanager.NewFromConfig(cfg)

	if secretArn, err = pkg.CreateSecret(context.TODO(), conn, secretName, value, "desc"); err != nil {
		panic("Couldn't create secret!: " + err.Error())
	}
	fmt.Printf("Created the arn %v\n", secretArn)

	if value, err = pkg.GetSecret(context.TODO(), conn, secretArn); err != nil {
		panic("Couldn't get secret value!")
	}
	fmt.Printf("it has the value \"%v\"\n", value)

	// You can also get the secret by name
	if value, err = pkg.GetSecret(context.TODO(), conn, secretName); err != nil {
		panic("Couldn't get secret value!")
	}
	fmt.Printf("it has the value \"%v\"\n", value)

	if err = pkg.UpdateSecret(context.TODO(), conn, secretArn, "correct horse battery staple"); err != nil {
		panic("Couldn't update secret!")
	}
	fmt.Println("The secret has been updated.")

	if value, err = pkg.GetSecret(context.TODO(), conn, secretArn); err != nil {
		panic("Couldn't get secret value!")
	}
	fmt.Printf("it has the value \"%v\"\n", value)

	var secretIds []string

	if secretIds, err = pkg.ListSecrets(context.TODO(), conn); err != nil {
		panic("Couldn't list secrets!")
	}

	fmt.Printf("There are %v secrets -- here's their IDs: \n", len(secretIds))
	for _, id := range secretIds {
		fmt.Println(id)
	}

	if err = pkg.DeleteSecret(context.TODO(), conn, secretArn); err != nil {
		panic("Couldn't delete secret!")
	}
	fmt.Printf("Deleted the secret with arn %v\n", secretArn)
}
