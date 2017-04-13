package ddb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/urfave/cli"
)

func ScanTable(ctx *cli.Context) error {
	tableName := ctx.String("table")
	printAsJson := ctx.Bool("json-output")

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
		if printAsJson {
			printItemAsJson(item)
		} else {
			printItem(item, kd)
			fmt.Println()
		}
	}

	return nil
}

func printItemAsJson(item map[string]*dynamodb.AttributeValue) {
	values := map[string]interface{}{}

	for k, v := range item {
		if v.S != nil {
			values[k] = *v.S
		} else if v.N != nil {
			c, err := strconv.Atoi(*v.N)
			if err != nil {
				log.Fatalln("ERROR", err)
			}

			values[k] = c
		}
	}

	err := json.NewEncoder(os.Stdout).Encode(values)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
}
