package cf

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/sam701/tcolor"
	"github.com/urfave/cli"
)

func printStackEventsCommand(ctx *cli.Context) error {
	stackName := ctx.Args().First()
	if stackName == "" {
		cli.ShowSubcommandHelp(ctx)
	} else {
		printStackEvents(stackName)
	}
	return nil
}

func printStackEvents(stackName string) {
	seenCount := 0
	readingEvents := true
	for {
		if !readingEvents {
			break
		}
		events := readStackEvents(stackName)
		if len(events) == 0 {
			break
		}

		if len(events) > seenCount {
			for _, event := range events[seenCount:] {
				fmt.Printf("%s %s %s\n%s%-40s %s\n",
					tcolor.Colorize(event.Timestamp.Format(time.RFC3339), tcolor.New().Foreground(tcolor.BrightWhite)),
					statusColor(fmt.Sprintf("%-20s", *event.ResourceStatus)),
					awsToString(event.ResourceType),
					strings.Repeat(" ", 21),
					awsToString(event.LogicalResourceId),
					awsToString(event.ResourceStatusReason),
				)
			}

			seenCount += len(events)
		}
		time.Sleep(5 * time.Second)
	}
}

func statusColor(status string) string {
	var c = tcolor.White
	if strings.Contains(status, "IN_PROGRESS") {
		c = tcolor.BrightYellow
	} else if strings.Contains(status, "_COMPLETE") {
		c = tcolor.BrightGreen
	} else if strings.Contains(status, "_DELETED") {
		c = tcolor.BrightYellow
	} else if strings.Contains(status, "_FAILED") {
		c = tcolor.BrightRed
	}

	return tcolor.Colorize(status, tcolor.New().Foreground(c))
}

var haveSuccessfulEventRetrieval = false

func readStackEvents(stackName string) []*cloudformation.StackEvent {
	var token *string = nil
	events := make([]*cloudformation.StackEvent, 0, 32)
	for {
		out, err := cfClient.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
			NextToken: token,
			StackName: aws.String(stackName),
		})
		if err != nil {
			if haveSuccessfulEventRetrieval {
				break
			} else {
				log.Fatalln("ERROR", err)
			}
		}
		haveSuccessfulEventRetrieval = true

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
