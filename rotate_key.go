package main

import (
	"fmt"
	"log"

	"github.com/sam701/awstools/config"
	"github.com/sam701/awstools/cred"
	"github.com/sam701/awstools/sess"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/urfave/cli"
)

func rotateMainAccountKeyAction(*cli.Context) error {
	client := iam.New(sess.New(config.Current.Profiles.MainAccount))
	rotateMainAccountKey(client)
	return nil
}

func rotateMainAccountKey(client *iam.IAM) {
	key, err := client.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(getUserName()),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	currentAccessKeyId := cred.GetMainAccountKeyId(config.Current.Profiles.MainAccount)
	cred.SaveCredentials(config.Current.Profiles.MainAccount, *key.AccessKey.AccessKeyId, *key.AccessKey.SecretAccessKey, "")
	fmt.Println("Created new access key")

	_, err = client.DeleteAccessKey(&iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(currentAccessKeyId),
		UserName:    aws.String(getUserName()),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	fmt.Println("Deleted old access key")
}
