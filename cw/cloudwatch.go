package cw

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/sam701/awstools/colors"
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

func parseTime(str string, layouts ...string) time.Time {
	for _, layout := range layouts {
		t, err := time.Parse(layout, str)
		if err == nil {
			n := time.Now()
			if t.Year() == 0 {
				t = time.Date(n.Year(), n.Month(), n.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), n.Location())
			} else {
				t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), n.Location())
			}
			return t
		}
	}
	log.Fatalln("Cannot parse time:", str)
	return time.Now()
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
	default:
		return parseTime(str, "15:04:05", "15:04", "2006-01-02", "2006-01-02T15:04", "2006-01-02T15:04:05")
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
	fmt.Printf("Searching for '%s' in %s in interval [%s, %s]\n", filter.pattern, groupName, filter.start, filter.end)
	params := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(groupName),
		StartTime:    aws.Int64(filter.start.Unix() * 1000),
		EndTime:      aws.Int64(filter.end.Unix() * 1000),
	}
	if filter.pattern != "" {
		params.FilterPattern = aws.String(filter.pattern)
	}
	for {
		err := client.FilterLogEventsPages(params, func(out *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
			fmt.Print(".")
			params.NextToken = out.NextToken

			for _, event := range out.Events {
				msg := *event.Message
				if filter.pattern != "" {
					msg = strings.Replace(msg, filter.pattern, colors.Match(filter.pattern), -1)
				}

				fmt.Printf("%s %s",
					colors.Timestamp(toTimeString(event.Timestamp)),
					msg)
			}
			return true
		})
		if err == nil {
			break
		} else if e, ok := err.(awserr.Error); ok {
			if e.Code() != "ThrottlingException" {
				log.Fatalln("ERROR", err)
			}
		}
	}
}

func listGroups() {
	client.DescribeLogGroupsPages(&cloudwatchlogs.DescribeLogGroupsInput{}, func(out *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
		for _, v := range out.LogGroups {
			fmt.Println(*v.LogGroupName)
		}
		return true
	})
}
