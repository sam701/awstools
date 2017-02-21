package ddb

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/urfave/cli"
)

func ScanTable(ctx *cli.Context) error {
	tableName := ctx.String("table")

	if tableName == "" {
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	c := client()
	out, err := c.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	kd := describeTableKey(tableName, c)

	limit := ctx.Int("row-limit")
	for i, item := range out.Items {
		if i >= limit {
			break
		}
		printItem(item, kd)
		fmt.Println()
	}

	return nil
}
