package cf

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/sam701/awstools/sess"
	"github.com/urfave/cli"
)

func deleteStack(ctx *cli.Context) error {
	name := ctx.Args().First()
	if name == "" {
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	cfClient = cloudformation.New(sess.FromEnvVar())
	_, err := cfClient.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String(name),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	printStackEvents(name)
	return nil
}
