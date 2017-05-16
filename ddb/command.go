package ddb

import (
	"github.com/urfave/cli"
)

func Command() cli.Command {
	return cli.Command{
		Name:      "dynamodb",
		ShortName: "ddb",
		Usage:     "dynamodb commands",
		Subcommands: []cli.Command{
			{
				Name:   "list",
				Usage:  "list tables",
				Action: List,
			},
			{
				Name:      "describe",
				ShortName: "desc",
				Usage:     "describe table",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "table, t",
						Usage: "table name",
					},
				},
				Action: DescribeTable,
			},
			{
				Name:  "get",
				Usage: "get item",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "table, t",
						Usage: "table name",
					},
					cli.StringFlag{
						Name:  "hash-key, k",
						Usage: "hash key of the item",
					},
					cli.StringFlag{
						Name:  "range-key, r",
						Usage: "range key of the item",
					},
					cli.BoolFlag{
						Name:  "json-output, j",
						Usage: "output lines as JSON",
					},
				},
				Action: GetItem,
			},
			{
				Name:  "put",
				Usage: "read JSON objects line by line from stdin and put them into DDB",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "table, t",
						Usage: "table name",
					},
				},
				Action: putItem,
			},
			{
				Name:  "delete",
				Usage: "delete item",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "table, t",
						Usage: "table name",
					},
					cli.StringFlag{
						Name:  "hash-key, k",
						Usage: "hash key of the item",
					},
					cli.StringFlag{
						Name:  "range-key, r",
						Usage: "range key of the item",
					},
				},
				Action: deleteItemAction,
			},
			{
				Name:  "scan",
				Usage: "return first rows",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "table, t",
						Usage: "table name",
					},
					cli.IntFlag{
						Name:  "row-limit, l",
						Usage: "maximum rows to print",
						Value: 20,
					},
					cli.BoolFlag{
						Name:  "json-output, j",
						Usage: "output lines as JSON",
					},
				},
				Action: ScanTable,
			},
		},
	}

}
