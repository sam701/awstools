package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aybabtme/rgbterm"
	"github.com/codegangsta/cli"
)

func kinesisPrintRecords(c *cli.Context) error {
	client := kinesis.New(currentEnvVarSession())

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

	for _, shard := range dsr.StreamDescription.Shards {
		go searchInShard(client, streamName, shard.ShardId, c.StringSlice("pattern"))
	}

	ch := make(chan bool)
	<-ch

	return nil
}

func searchInShard(client *kinesis.Kinesis, streamName string, shardId *string, patterns []string) {
	itOut, err := client.GetShardIterator(&kinesis.GetShardIteratorInput{
		StreamName:        aws.String(streamName),
		ShardId:           shardId,
		ShardIteratorType: aws.String("LATEST"),
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
				printKinesisEvent(record, patterns)
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
					printKinesisEvent(record, patterns)
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func printKinesisEvent(record *kinesis.Record, patterns []string) {
	tt := record.ApproximateArrivalTimestamp.Format(time.RFC3339)
	tt = rgbterm.FgString(tt, 100, 255, 100)
	fmt.Print(tt + " ")
	str := strings.TrimSpace(string(record.Data))
	if len(patterns) == 0 {
		fmt.Println(str)
	} else {
		for _, pat := range patterns {
			str = strings.Replace(str, pat, rgbterm.FgString(pat, 255, 100, 100), -1)
		}
		fmt.Println(str)
	}
}
