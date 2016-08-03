package cw

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aybabtme/rgbterm"
	"github.com/sam701/awstools/sess"
	"github.com/urfave/cli"
)

var client *cloudwatchlogs.CloudWatchLogs

func CloudwatchAction(c *cli.Context) error {
	client = cloudwatchlogs.New(sess.FromEnvVar())
	if c.Bool("list-groups") {
		listGroups()
	} else if group := c.String("group"); group != "" {
		grabInGroup(group, &filter{
			pattern: c.String("pattern"),
			start:   parseTimeBoundary(c.String("start")),
			end:     parseTimeBoundary(c.String("end")),
		})

	} else {
		cli.ShowCommandHelp(c, "cloudwatch")
	}
	return nil
}

type filter struct {
	pattern string
	start   time.Time
	end     time.Time
}

func parseTimeBoundary(str string) time.Time {
	now := time.Now()
	if str == "now" {
		return now
	}

	lastRune := str[len(str)-1]
	duration := time.Second
	switch lastRune {
	case 'h':
		duration = time.Hour
	case 'm':
		duration = time.Minute
	case 's':
		duration = time.Second
	}

	i, err := strconv.Atoi(str[:len(str)-1])
	if err != nil {
		log.Fatalln("Cannot parse", str, "Samples: -2h, -3m, -30s")
	}

	return now.Add(time.Duration(i * int(duration)))
}

func toTimeString(millis *int64) string {
	if millis == nil {
		return ""
	}
	return time.Unix(*millis/1000, 0).Format(time.RFC3339)
}

func grabInGroup(groupName string, filter *filter) {
	params := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(groupName),
		StartTime:    aws.Int64(filter.start.Unix() * 1000),
		EndTime:      aws.Int64(filter.end.Unix() * 1000),
	}
	if filter.pattern != "" {
		params.FilterPattern = aws.String(filter.pattern)
	}
	err := client.FilterLogEventsPages(params, func(out *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		fmt.Print(".")
		for _, event := range out.Events {
			msg := *event.Message
			if filter.pattern != "" {
				msg = strings.Replace(msg, filter.pattern, rgbterm.FgString(filter.pattern, 255, 100, 100), -1)
			}

			fmt.Printf("%s %s",
				rgbterm.FgString(toTimeString(event.Timestamp), 130, 255, 130),
				msg)
		}
		return true
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	fmt.Println("return")

}

func listGroups() {
	client.DescribeLogGroupsPages(&cloudwatchlogs.DescribeLogGroupsInput{}, func(out *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
		for _, v := range out.LogGroups {
			fmt.Println(*v.LogGroupName)
		}
		return true
	})
}
