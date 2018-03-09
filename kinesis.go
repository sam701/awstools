package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
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
	var atTimestamp *time.Time
	if str := c.String("at-timestamp"); str != "" {
		t, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return errors.New("Cannot parse time " + err.Error())
		}
		atTimestamp = &t
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
			atTimestamp,
			!c.Bool("no-timestamp"),
			c.Bool("gunzip"),
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
	atTimestamp *time.Time,
	printTimestamp bool,
	gunzip bool,
	eventsToPrint chan<- *eventToPrint) {

	siInput := &kinesis.GetShardIteratorInput{
		StreamName: aws.String(streamName),
		ShardId:    shardId,
	}
	siInput.ShardIteratorType = aws.String("LATEST")
	if trimHorizon {
		siInput.ShardIteratorType = aws.String("TRIM_HORIZON")
	} else if atTimestamp != nil {
		siInput.ShardIteratorType = aws.String("AT_TIMESTAMP")
		siInput.Timestamp = atTimestamp
	}
	itOut, err := client.GetShardIterator(siInput)
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
			if gunzip {
				r, err := gzip.NewReader(bytes.NewReader(record.Data))
				if err != nil {
					log.Println("Cannot unzip data", err, string(record.Data))
					continue
				}

				var output bytes.Buffer
				_, err = io.Copy(&output, r)
				if err != nil {
					log.Println("Cannot unzip data", err, string(record.Data))
					continue
				}
				r.Close()
				record.Data = output.Bytes()
			}
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
