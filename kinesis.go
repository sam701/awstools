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

		itOut, err := client.GetShardIterator(&kinesis.GetShardIteratorInput{
			StreamName:        aws.String(streamName),
			ShardId:           shard.ShardId,
			ShardIteratorType: aws.String("LATEST"),
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		go searchInShard(client, itOut.ShardIterator, c.StringSlice("pattern"))
	}

	ch := make(chan bool)
	<-ch

	return nil
}

func searchInShard(client *kinesis.Kinesis, shardIterator *string, patterns []string) {
	for {
		rOut, err := client.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		for _, record := range rOut.Records {
			if len(patterns) == 0 {
				printKinesisEvent(record, "")
			} else {
				str := string(record.Data)
				for _, pat := range patterns {
					if strings.Contains(str, pat) {
						printKinesisEvent(record, pat)
					}
				}
			}
		}
	}
}

func printKinesisEvent(record *kinesis.Record, pattern string) {
	tt := record.ApproximateArrivalTimestamp.Format(time.RFC3339)
	tt = rgbterm.FgString(tt, 100, 255, 100)
	fmt.Print(tt + " ")
	str := strings.TrimSpace(string(record.Data))
	if pattern == "" {
		fmt.Println(str)
	} else {
		str = strings.Replace(str, pattern, rgbterm.FgString(pattern, 255, 100, 100), -1)
		fmt.Println(str)
	}
}
