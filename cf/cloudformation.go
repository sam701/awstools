package cf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/sam701/awstools/sess"
	"github.com/urfave/cli"
)

var cfClient *cloudformation.CloudFormation

func HandleCloudformation(c *cli.Context) error {
	cfClient = cloudformation.New(sess.FromEnvVar())
	if searchPattern := c.String("search"); searchPattern != "" {
		searchStack(searchPattern)
	} else if stackToDelete := c.String("delete"); stackToDelete != "" {
		deleteStack(stackToDelete)
	} else if eventsStackName := c.String("events"); eventsStackName != "" {
		printStackEvents(aws.String(eventsStackName))
	} else {
		cli.ShowCommandHelp(c, "cloudformation")
	}
	return nil
}

func awsToString(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}
