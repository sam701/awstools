package main

import (
	"fmt"
	"log"

	"github.com/sam701/awstools/cred"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/codegangsta/cli"
)

func rotateBastionKey(*cli.Context) error {
	client := iam.New(newSession(theConfig.Profiles.Bastion))
	key, err := client.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(getUserName()),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	currentAccessKeyId := cred.GetBastionKeyId(theConfig.Profiles.Bastion)
	cred.SaveCredentials(theConfig.Profiles.Bastion, *key.AccessKey.AccessKeyId, *key.AccessKey.SecretAccessKey, "")
	fmt.Println("Created new access key")

	_, err = client.DeleteAccessKey(&iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(currentAccessKeyId),
		UserName:    aws.String(getUserName()),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	fmt.Println("Deleted old access key")

	return nil
}
