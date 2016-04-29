package main

import (
	"fmt"
	"log"
	"sam701/awstools/cred"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/codegangsta/cli"
)

func rotateBastionKey(*cli.Context) {
	client := iam.New(newSession(theConfig.Bastion))
	key, err := client.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(getUserName()),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	currentAccessKeyId := cred.GetBastionKeyId(theConfig.Bastion)
	cred.SaveCredentials(theConfig.Bastion, *key.AccessKey.AccessKeyId, *key.AccessKey.SecretAccessKey, "")
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
