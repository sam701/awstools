package main

import (
	"fmt"
	"os"

	"github.com/sam701/awstools/assume"
	"github.com/sam701/awstools/cf"
	"github.com/sam701/awstools/config"
	"github.com/sam701/awstools/cw"
	"github.com/sam701/awstools/ddb"
	"github.com/sam701/tcolor"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "awstools"
	app.Version = "0.13.0"
	app.Usage = "AWS tools"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to config.toml file (default: ~/.config/awstools/config.toml)",
		},
		cli.BoolFlag{
			Name:  "no-color",
			Usage: "turn off color output",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "assume",
			Usage:     "assume role on a specified account",
			ArgsUsage: "<account name> <role to assume>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "export, e",
					Usage: "export shell script to `PATH`",
				},
				cli.BoolFlag{
					Name:  "export-profile, p",
					Usage: "export AWS_PROFILE instead of AWS_SECRET_ACCESS_KEY and AWS_SESSION_TOKEN",
				},
			},
			Action: assume.AssumeRole,
		},
		{
			Name:   "accounts",
			Usage:  "print known accounts",
			Action: actionPrintKnownAccounts,
		},
		{
			Name:      "ec2",
			Usage:     "print EC2 instances and ELBs",
			ArgsUsage: "<EC2 instance tag substring>",
			Action:    actionDescribeEC2,
		},
		cf.Command(),
		{
			Name:      "rotate-main-account-key",
			ShortName: "r",
			Usage:     "create a new access key for main account and delete the current one",
			Action:    assume.RotateMainAccountKeyAction,
		},
		ddb.Command(),
		{
			Name:      "kms",
			Usage:     "encrypt/decrypt text",
			ArgsUsage: "<text to de-/encrypt>",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "list-keys, l",
					Usage: "list KMS keys",
				},
				cli.BoolFlag{
					Name:  "decrypt, d",
					Usage: "decrypt base64 encoded string",
				},
				cli.BoolFlag{
					Name:  "encrypt, e",
					Usage: "encrypt and base64 encode string",
				},
				cli.StringFlag{
					Name:  "key-id",
					Usage: "use `KEYID` for encryption/decryption",
				},
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "print only de-/encrypted text",
				},
			},
			Action: kmsAction,
		},
		{
			Name:  "kinesis",
			Usage: "print records from kinesis streams",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "list-streams, l",
					Usage: "list kinesis streams",
				},
				cli.StringFlag{
					Name:  "search-stream, s",
					Usage: "search `STREAM`",
				},
				cli.BoolFlag{
					Name:  "trim-horizon",
					Usage: "Search from the last untrimmed record in the shards",
				},
				cli.StringSliceFlag{
					Name:  "pattern, p",
					Usage: "search for (case sensitive) `PATTERN`",
				},
				cli.BoolFlag{
					Name:  "no-timestamp",
					Usage: "do not print record timestamps",
				},
			},
			Action: kinesisPrintRecords,
		},
		{
			Name:      "cloudwatch",
			ShortName: "cw",
			Usage:     "search in cloudwatch logs",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "list-groups, l",
					Usage: "list CloudWatch groups",
				},
				cli.StringFlag{
					Name:  "group, g",
					Usage: "`GROUP` to grab (or a unique substring)",
				},
				cli.StringFlag{
					Name:  "pattern, p",
					Usage: "`PATTERN` to grab",
				},
				cli.StringFlag{
					Name:  "start",
					Usage: "start `TIME` to grab",
					Value: "-24h",
				},
				cli.StringFlag{
					Name:  "end",
					Usage: "end `TIME` to grab",
					Value: "now",
				},
				cli.IntFlag{
					Name:  "duration",
					Usage: "duration in seconds",
				},
			},
			Action: cw.CloudwatchAction,
		},
	}
	app.Before = func(c *cli.Context) error {
		config.Read(c.String("config"))

		if c.GlobalBool("no-color") {
			tcolor.ColorOn = false
		}
		return nil
	}
	app.Run(os.Args)
}

func actionPrintKnownAccounts(c *cli.Context) error {
	for name, accountId := range config.Current.Accounts {
		fmt.Println(name, "=", accountId)
	}
	return nil
}
