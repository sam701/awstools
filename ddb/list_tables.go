package ddb

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sam701/awstools/sess"
	"github.com/sam701/tcolor"
	"github.com/urfave/cli"
)

func client() *dynamodb.DynamoDB {
	return dynamodb.New(sess.FromEnvVar())
}

func List(ctx *cli.Context) error {
	c := client()
	out, err := c.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	color := tcolor.New().Foreground(tcolor.BrightYellow)
	for _, table := range out.TableNames {
		fmt.Println(tcolor.Colorize(*table, color))
	}
	return nil
}
