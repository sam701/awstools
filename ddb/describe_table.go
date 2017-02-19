package ddb

import (
	"fmt"
	"log"
	"strconv"

	"time"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	humanize "github.com/dustin/go-humanize"
	"github.com/sam701/awstools/colors"
	"github.com/sam701/tcolor"
	"github.com/urfave/cli"
)

func DescribeTable(ctx *cli.Context) error {
	tableName := ctx.String("table")
	if tableName == "" {
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	c := client()
	out, err := c.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	desc := out.Table
	fmt.Printf("%s\n",
		tcolor.Colorize(*desc.TableName, tcolor.New().Foreground(tcolor.BrightYellow).Bold()),
	)

	printProperties(
		"Status", *desc.TableStatus,
		"Creation Time", desc.CreationDateTime.Format(time.RFC3339),
		"Item Count", humanize.Comma(int64(*desc.ItemCount)),
		"Database Size", humanize.Bytes(uint64(*desc.TableSizeBytes)),
		"Key Schema", keySchemaString(desc),
		"Provisioned Read", humanize.Comma(*desc.ProvisionedThroughput.ReadCapacityUnits),
		"Provisioned Write", humanize.Comma(*desc.ProvisionedThroughput.WriteCapacityUnits),
		"Last Increase", desc.ProvisionedThroughput.LastIncreaseDateTime.Format(time.RFC3339),
		"Last Decrease", desc.ProvisionedThroughput.LastDecreaseDateTime.Format(time.RFC3339),
		"Decreases Today", strconv.Itoa(int(*desc.ProvisionedThroughput.NumberOfDecreasesToday)),
	)

	return nil
}

func keySchemaString(desc *dynamodb.TableDescription) string {
	parts := []string{}
	for i, el := range desc.KeySchema {
		parts = append(parts, *el.KeyType+": "+*el.AttributeName+"("+*desc.AttributeDefinitions[i].AttributeType+")")
	}
	return strings.Join(parts, ", ")
}

func printProperties(pairs ...string) {
	for i := 0; i < len(pairs); i += 2 {
		prop := pairs[i]
		value := pairs[i+1]

		fmt.Printf("  %s %s\n",
			colors.Property(fmt.Sprintf("%-20s", prop+":")),
			value)

	}
}
