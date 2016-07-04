package cf

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aybabtme/rgbterm"
)

func printStackEvents(stackName *string) {
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
				fmt.Printf("%s %-20s %s\n    %-40s %s\n",
					rgbterm.FgString(event.Timestamp.Format(time.RFC3339), 130, 255, 130),
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
}

func readStackEvents(stackName *string) []*cloudformation.StackEvent {
	var token *string = nil
	events := make([]*cloudformation.StackEvent, 0, 32)
	for {
		out, err := cfClient.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
			NextToken: token,
			StackName: stackName,
		})
		if err != nil {
			break
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
