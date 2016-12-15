package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/sam701/awstools/colors"
	"github.com/sam701/awstools/sess"
	"github.com/urfave/cli"
)

func kinesisPrintRecords(c *cli.Context) error {
	client := kinesis.New(sess.FromEnvVar())

	if c.Bool("list-streams") {
		out, err := client.ListStreams(&kinesis.ListStreamsInput{})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		for _, s := range out.StreamNames {
			fmt.Println(*s)
		}
		return nil
	}
	streamName := c.String("search-stream")
	if streamName == "" {
		cli.ShowCommandHelp(c, "kinesis")
		return nil
	}

	dsr, err := client.DescribeStream(&kinesis.DescribeStreamInput{
		StreamName: aws.String(streamName),
		Limit:      aws.Int64(1000),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	eventsChan := make(chan *eventToPrint)
	go printEvents(eventsChan)

	for _, shard := range dsr.StreamDescription.Shards {
		go searchInShard(client, streamName, shard.ShardId,
			c.StringSlice("pattern"),
			c.Bool("trim-horizon"),
			!c.Bool("no-timestamp"),
			eventsChan,
		)
	}

	ch := make(chan bool)
	<-ch

	return nil
}

func searchInShard(
	client *kinesis.Kinesis,
	streamName string,
	shardId *string,
	patterns []string,
	trimHorizon bool,
	printTimestamp bool,
	eventsToPrint chan<- *eventToPrint) {

	shardIteratorType := "LATEST"
	if trimHorizon {
		shardIteratorType = "TRIM_HORIZON"
	}
	itOut, err := client.GetShardIterator(&kinesis.GetShardIteratorInput{
		StreamName:        aws.String(streamName),
		ShardId:           shardId,
		ShardIteratorType: aws.String(shardIteratorType),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	shardIterator := itOut.ShardIterator

	for {
		if shardIterator == nil {
			return
		}
		rOut, err := client.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		shardIterator = rOut.NextShardIterator

		for _, record := range rOut.Records {
			if len(patterns) == 0 {
				eventsToPrint <- &eventToPrint{record, patterns, printTimestamp}
			} else {
				str := string(record.Data)
				matched := true
				for _, pat := range patterns {
					if !strings.Contains(str, pat) {
						matched = false
						break
					}
				}
				if matched {
					eventsToPrint <- &eventToPrint{record, patterns, printTimestamp}
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}

type eventToPrint struct {
	record         *kinesis.Record
	patterns       []string
	printTimestamp bool
}

func printEvents(events <-chan *eventToPrint) {
	for event := range events {

		if event.printTimestamp {
			tt := event.record.ApproximateArrivalTimestamp.Format(time.RFC3339)
			tt = colors.Timestamp(tt)
			fmt.Print(tt + " ")
		}
		str := strings.TrimSpace(string(event.record.Data))
		if len(event.patterns) == 0 {
			fmt.Println(str)
		} else {
			for _, pat := range event.patterns {
				str = strings.Replace(str, pat, colors.Match(pat), -1)
			}
			fmt.Println(str)
		}
	}
}
