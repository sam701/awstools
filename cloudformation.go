package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aybabtme/rgbterm"
	"github.com/codegangsta/cli"
)

func printStacks(c *cli.Context) error {
	if searchPattern := c.String("search"); searchPattern != "" {
		cf := cloudformation.New(currentEnvVarSession())
		out, err := cf.DescribeStacks(&cloudformation.DescribeStacksInput{})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		for _, stack := range out.Stacks {
			if strings.Contains(*stack.StackName, searchPattern) {
				formatStack(stack)
			}
		}
	} else if stackToDelete := c.String("delete"); stackToDelete != "" {
		deleteStack(stackToDelete)
	} else {
		cli.ShowCommandHelp(c, "cloudformation")
	}
	return nil
}

func deleteStack(name string) bool {
	cf := cloudformation.New(currentEnvVarSession())
	awsName := aws.String(name)
	_, err := cf.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: awsName,
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	seenCount := 0
	readingEvents := true
	for {
		if !readingEvents {
			break
		}
		events := readStackEvents(cf, awsName)

		if len(events) > seenCount {
			for _, event := range events[seenCount:] {
				fmt.Printf("%s %-20s %-30s %-40s %s\n",
					event.Timestamp.Format(time.RFC3339),
					awsToString(event.ResourceStatus),
					awsToString(event.ResourceType),
					awsToString(event.LogicalResourceId),
					awsToString(event.ResourceStatusReason),
				)
			}

			seenCount += len(events)
		}
		time.Sleep(5 * time.Second)
	}

	return true
}

func awsToString(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}

func readStackEvents(cf *cloudformation.CloudFormation, stackName *string) []*cloudformation.StackEvent {
	var token *string = nil
	events := make([]*cloudformation.StackEvent, 0, 32)
	for {
		out, err := cf.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
			NextToken: token,
			StackName: stackName,
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		events = append(events, out.StackEvents...)

		token = out.NextToken
		if token == nil {
			break
		}
	}

	// reverse
	for i, j := 0, len(events)-1; i < j; i, j = i+1, j-1 {
		events[i], events[j] = events[j], events[i]
	}

	return events
}

func formatStack(stack *cloudformation.Stack) {
	fmt.Printf("%s %s\n",
		rgbterm.FgString(fmt.Sprintf("%-30s", *stack.StackName), 255, 255, 255),
		rgbterm.FgString(*stack.StackStatus, 130, 255, 130))
	fmt.Println("  Parameters:")
	for _, par := range stack.Parameters {
		fmt.Printf("    %-35s %s\n", *par.ParameterKey, *par.ParameterValue)
	}
	fmt.Println("  Outputs:")
	for _, out := range stack.Outputs {
		fmt.Printf("    %-35s %s\n", *out.OutputKey, *out.OutputValue)
	}
	fmt.Println()
}
