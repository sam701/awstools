package ddb

import (
	"fmt"
	"log"

	"strings"

	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sam701/awstools/colors"
	"github.com/sam701/tcolor"
	"github.com/urfave/cli"
)

func GetItem(ctx *cli.Context) error {
	table := ctx.String("table")
	hashKey := ctx.String("hash-key")
	rangeKey := ctx.String("range-key")

	if table == "" || hashKey == "" {
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	c := client()

	kd := describeTableKey(table, c)

	out, err := c.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       createKey(kd, hashKey, rangeKey),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	printItem(out.Item, kd)

	return nil
}

type fieldValue struct {
	value     string
	valueType string
}

type itemType map[string]*dynamodb.AttributeValue

func (i itemType) getStringValue(key string) *fieldValue {
	val := i[key]
	if s := val.S; s != nil {
		return &fieldValue{*s, "S"}
	} else if n := val.N; n != nil {
		return &fieldValue{*n, "N"}
	} else {
		log.Fatalln("Cannot interpret type of field", key)
	}

	return nil
}

func printItem(item itemType, kd *keyDesc) {
	if len(item) == 0 {
		fmt.Println("No such item")
	}

	maxFieldNameSize := 0
	keys := []string{}

	for k, _ := range item {
		keys = append(keys, k)
		if len(k) > maxFieldNameSize {
			maxFieldNameSize = len(k)
		}
	}

	sort.Strings(keys)

	printOneProp(kd.hashKey.keyName, "HASH", maxFieldNameSize, item.getStringValue(kd.hashKey.keyName))
	if kd.rangeKey != nil {
		printOneProp(kd.rangeKey.keyName, "RANGE", maxFieldNameSize, item.getStringValue(kd.rangeKey.keyName))
	}

	for _, k := range keys {
		if kd.hashKey.keyName == k || kd.rangeKey != nil && kd.rangeKey.keyName == k {
			continue
		}
		printOneProp(k, "prop", maxFieldNameSize, item.getStringValue(k))
	}
}

var (
	colorRange  = tcolor.New().ForegroundGray(18).Italic().Bold()
	colorHash   = colorRange.Underline()
	colorString = tcolor.New().Foreground(tcolor.BrightBlue)
	colorNumber = tcolor.New().Foreground(tcolor.BrightRed)
)

func printOneProp(key, keyType string, alignAtWidth int, value *fieldValue) {
	var col tcolor.Color
	switch keyType {
	case "prop":
		col = colors.PropertyColor
	case "HASH":
		col = colorHash
	case "RANGE":
		col = colorRange
	}
	fmt.Print(tcolor.Colorize(key+":", col), strings.Repeat(" ", alignAtWidth-len(key)+3))

	col = ""
	switch value.valueType {
	case "N":
		col = colorNumber
	case "S":
		col = colorString
	}
	fmt.Println(tcolor.Colorize(value.value, col))
}

func createAttributeValue(value string, valueType string) *dynamodb.AttributeValue {
	av := &dynamodb.AttributeValue{}
	val := aws.String(value)
	switch valueType {
	case "S":
		av.S = val
	case "N":
		av.N = val
	}
	return av
}

func createKey(kd *keyDesc, hashValue, rangeValue string) itemType {
	key := itemType{}
	key[kd.hashKey.keyName] = createAttributeValue(hashValue, kd.hashKey.keyValueType)

	if kd.rangeKey != nil {
		key[kd.rangeKey.keyName] = createAttributeValue(rangeValue, kd.rangeKey.keyValueType)
	}

	return key
}

type fieldDesc struct {
	keyName      string
	keyValueType string
}

type keyDesc struct {
	hashKey  *fieldDesc
	rangeKey *fieldDesc
}

func describeTableKey(tableName string, c *dynamodb.DynamoDB) *keyDesc {
	out, err := c.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	kd := &keyDesc{
		hashKey: &fieldDesc{},
	}

	for i, sch := range out.Table.KeySchema {
		if *sch.KeyType == "HASH" {
			kd.hashKey.keyName = *sch.AttributeName
			kd.hashKey.keyValueType = *out.Table.AttributeDefinitions[i].AttributeType
		} else {
			kd.rangeKey = &fieldDesc{
				keyName:      *sch.AttributeName,
				keyValueType: *out.Table.AttributeDefinitions[i].AttributeType,
			}
		}
	}

	return kd
}
