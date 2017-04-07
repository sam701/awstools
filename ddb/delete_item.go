package ddb

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/urfave/cli"
)

func deleteItemAction(ctx *cli.Context) error {
	table := ctx.String("table")
	hashKey := ctx.String("hash-key")
	rangeKey := ctx.String("range-key")

	if table == "" || hashKey == "" {
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	c := client()

	kd := describeTableKey(table, c)
	if kd.rangeKey != nil && rangeKey == "" {
		return errors.New("Missing range-key")
	}

	key := map[string]*dynamodb.AttributeValue{}
	key[kd.hashKey.keyName] = createAttributeValue(hashKey, kd.hashKey.keyValueType)
	if kd.rangeKey != nil {
		key[kd.rangeKey.keyName] = createAttributeValue(rangeKey, kd.rangeKey.keyValueType)
	}

	_, err := c.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key:       key,
	})

	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return nil
}
