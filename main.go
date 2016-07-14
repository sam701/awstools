package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/sam701/awstools/cf"
	"github.com/sam701/awstools/config"
)

func main() {
	app := cli.NewApp()
	app.Name = "awstools"
	app.Version = "0.8.0"
	app.Usage = "AWS tools"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to config.toml file (default: ~/.config/awstools/config.toml)",
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
					Usage: "path to export the shell script to source in",
				},
			},
			Action: actionAssumeRole,
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
		{
			Name:      "cloudformation",
			ShortName: "cf",
			Usage:     "print CloudFormation stacks information",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "search, s",
					Usage: "stack name substring",
				},
				cli.StringFlag{
					Name:  "events, e",
					Usage: "print stack's events",
				},
				cli.StringFlag{
					Name:  "delete",
					Usage: "delete stack",
				},
			},
			Action: cf.HandleCloudformation,
		},
		{
			Name:      "rotate-main-account-key",
			ShortName: "r",
			Usage:     "create a new access key for main account and delete the current one",
			Action:    rotateMainAccountKeyAction,
		},
		{
			Name:      "kms",
			Usage:     "encrypt/decrypt text",
			ArgsUsage: "<text to decrypt>",
			Flags: []cli.Flag{
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
					Usage: "",
				},
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "print only encrypted/decrypted text",
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
					Usage: "stream to search",
				},
				cli.StringSliceFlag{
					Name:  "pattern, p",
					Usage: "pattern to search for (case sensitive)",
				},
			},
			Action: kinesisPrintRecords,
		},
	}
	app.Before = func(c *cli.Context) error {
		config.Read(c.String("config"))
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
