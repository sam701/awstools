package ddb

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"errors"

	"strconv"

	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/urfave/cli"
)

func putItem(ctx *cli.Context) error {
	table := ctx.String("table")
	if table == "" {
		return errors.New("No table name was provided")
	}

	dec := json.NewDecoder(os.Stdin)
	c := client()
	for {
		var data map[string]interface{}
		err := dec.Decode(&data)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		}

		_, err = c.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item:      createDdbItem(data),
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
	}
	return nil
}

func createDdbItem(data map[string]interface{}) map[string]*dynamodb.AttributeValue {
	item := map[string]*dynamodb.AttributeValue{}
	for k, v := range data {
		item[k] = createAttributeValueFramInterface(v)
	}
	return item
}

func createAttributeValueFramInterface(value interface{}) *dynamodb.AttributeValue {
	switch val := value.(type) {
	case string:
		return createAttributeValue(val, "S")
	case float64:
		return createAttributeValue(strconv.FormatFloat(val, 'f', -1, 64), "N")
	default:
		log.Fatalln("Unsupported type", val, reflect.TypeOf(val))
	}
	return nil
}
