package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var theConfig *configuration

func main() {
	app := cli.NewApp()
	app.Name = "awstools"
	app.Version = "0.4.0"
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
			Aliases:   []string{"e"},
			Usage:     "print EC2 instances and ELBs for the given stack",
			ArgsUsage: "<stack name prefix>",
			Action:    actionDescribeEC2,
		},
		{
			Name:    "rotate-bastion-key",
			Aliases: []string{"r"},
			Usage:   "create a new bastion access key and delete the current one",
			Action:  rotateBastionKey,
		},
	}
	app.Before = func(c *cli.Context) error {
		theConfig = readConfig(c.String("config"))
		return nil
	}
	app.Run(os.Args)
}

func actionPrintKnownAccounts(c *cli.Context) error {
	for name, accountId := range theConfig.Accounts {
		fmt.Println(name, "=", accountId)
	}
	return nil
}
