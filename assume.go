package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/sam701/awstools/cred"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/codegangsta/cli"
)

var scriptOutput io.Writer = os.Stdout

func actionAssumeRole(c *cli.Context) {
	if len(c.Args()) == 2 {
		var account, role string
		account = c.Args().Get(0)
		role = c.Args().Get(1)
		exportScriptPath := c.String("export")
		if exportScriptPath != "" {
			f, err := os.OpenFile(exportScriptPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				log.Fatalln("ERROR:", err)
			}
			defer f.Close()
			scriptOutput = f
		}
		assumeRole(account, role)
	} else {
		cli.ShowCommandHelp(c, "assume")
	}
}

func assumeRole(account, role string) {
	account = adjustAccountName(account)
	role = adjustRoleName(role)

	err := tryToAssumeRole(account, role)
	if err != nil {
		getBastionSessionToken()
		err = tryToAssumeRole(account, role)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func adjustAccountName(account string) string {
	var candidate string
	for k, _ := range theConfig.Accounts {
		if strings.Contains(k, account) {
			if candidate == "" {
				candidate = k
			} else {
				log.Fatalf("Ambiguous account name. Possible matches: %s, %s\n", candidate, k)
			}
		}
	}
	if candidate == "" {
		log.Fatalln("No such account:", account)
	}
	return candidate
}

func adjustRoleName(role string) string {
	if role[0] == 'r' {
		return "ReadOnlyAccess"
	} else if role[0] == 'w' {
		return "PowerUserAccess"
	} else {
		return role
	}
}

var userName string

func getUserName() string {
	if userName == "" {
		client := iam.New(newSession(theConfig.Profiles.Bastion))
		data, err := client.GetUser(&iam.GetUserInput{})
		if err != nil {
			log.Fatalln("ERROR:", err)
		}
		userName = *data.User.UserName
	}
	return userName
}

func readMfaToken() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("MFA token: ")
	scanner.Scan()
	return scanner.Text()
}

func getBastionSessionToken() {
	token := readMfaToken()
	session := newSession(theConfig.Profiles.Bastion)
	stsClient := sts.New(session)
	data, err := stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		SerialNumber: aws.String(fmt.Sprintf("arn:aws:iam::%s:mfa/%s",
			accountId(theConfig.Profiles.Bastion),
			getUserName())),
		TokenCode: aws.String(token),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	persistSharedCredentials(data.Credentials, theConfig.Profiles.BastionMfa)
}

func newSession(name string) *session.Session {
	return session.New(&aws.Config{
		Region:      aws.String(theConfig.DefaultRegion),
		Credentials: credentials.NewSharedCredentials("", name),
	})
}

func accountId(accountName string) string {
	accountId := theConfig.Accounts[accountName]
	if accountId == "" {
		log.Fatalln("Unknown account name:", accountName)
	}
	return accountId
}

func tryToAssumeRole(accountName, role string) error {
	session := newSession(theConfig.Profiles.BastionMfa)
	accountId := accountId(accountName)

	stsClient := sts.New(session)
	assumeData, err := stsClient.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, role)),
		RoleSessionName: aws.String(getUserName()),
	})
	if err != nil {
		return err
	}

	profile := fmt.Sprintf("%s %s", accountName, role)
	persistSharedCredentials(assumeData.Credentials, profile)
	printExport(assumeData.Credentials)
	return nil
}

func printExport(cred *sts.Credentials) {
	shell := path.Base(os.Getenv("SHELL"))
	var pattern string
	if shell == "fish" {
		pattern = "set -xg %s \"%s\"\n"
	} else {
		pattern = "export %s=\"%s\"\n"
	}

	exp := func(key, value string) {
		fmt.Fprintf(scriptOutput, pattern, key, value)
	}

	exp("AWS_ACCESS_KEY_ID", *cred.AccessKeyId)
	exp("AWS_SECRET_ACCESS_KEY", *cred.SecretAccessKey)
	exp("AWS_SESSION_TOKEN", *cred.SessionToken)
}

func persistSharedCredentials(c *sts.Credentials, profile string) {
	cred.SaveCredentials(profile, *c.AccessKeyId, *c.SecretAccessKey, *c.SessionToken)
}
