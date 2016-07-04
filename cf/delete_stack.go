package cf

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func deleteStack(name string) bool {
	awsName := aws.String(name)
	_, err := cfClient.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: awsName,
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	printStackEvents(awsName)
	return true
}
