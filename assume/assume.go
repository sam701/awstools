package assume

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sam701/awstools/config"
	"github.com/sam701/awstools/cred"
	"github.com/sam701/awstools/sess"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/urfave/cli"
)

var exportProfile = false

func AssumeRole(c *cli.Context) error {
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
		exportProfile = c.Bool("export-profile")
		assumeRole(account, role)
	} else {
		cli.ShowCommandHelp(c, "assume")
	}
	return nil
}

func assumeRole(account, role string) {
	account = adjustAccountName(account)
	role = adjustRoleName(role)

	err := tryToAssumeRole(account, role)
	if err != nil {

		getMainAccountMfaSessionToken()
		err = tryToAssumeRole(account, role)
		if err != nil {
			log.Fatalln(err)
		}

		if needRotateKey() {
			rotateMainAccountKey(iam.New(sess.New(config.Current.Profiles.MainAccountMfaSession)))
		}
	}
}

func needRotateKey() bool {
	session := sess.New(config.Current.Profiles.MainAccountMfaSession)
	cl := iam.New(session)

	keyId := cred.GetMainAccountKeyId(config.Current.Profiles.MainAccount)
	out, err := cl.ListAccessKeys(&iam.ListAccessKeysInput{})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	var creationTime time.Time
	for _, md := range out.AccessKeyMetadata {
		if keyId == *md.AccessKeyId {
			creationTime = *md.CreateDate
			break
		}
	}
	if creationTime.IsZero() {
		log.Fatalln("Cannot get creation time for key", keyId)
	}

	return int(time.Now().Sub(creationTime).Minutes()) >= config.Current.KeyRotationIntervalMinutes
}

func adjustAccountName(account string) string {
	var candidate string
	for k, _ := range config.Current.Accounts {
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
		client := iam.New(sess.New(config.Current.Profiles.MainAccount))
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

func getMainAccountMfaSessionToken() {
	token := readMfaToken()
	session := sess.New(config.Current.Profiles.MainAccount)
	stsClient := sts.New(session)
	data, err := stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		SerialNumber: aws.String(fmt.Sprintf("arn:aws:iam::%s:mfa/%s",
			accountId(config.Current.Profiles.MainAccount),
			getUserName())),
		TokenCode: aws.String(token),
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	persistSharedCredentials(data.Credentials, config.Current.Profiles.MainAccountMfaSession)
}

func accountId(accountName string) string {
	accountId := config.Current.Accounts[accountName]
	if accountId == "" {
		log.Fatalln("Unknown account name:", accountName)
	}
	return accountId
}

func tryToAssumeRole(accountName, role string) error {
	session := sess.New(config.Current.Profiles.MainAccountMfaSession)
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
	if exportProfile {
		printExportProfile(profile)
	} else {
		printExportKeyAndToken(assumeData.Credentials)
	}

	cred.SetProfileRegion(profile, config.Current.DefaultRegion)

	if assumeData.Credentials.Expiration != nil {
		saveProfileExpirationTimestamp(profile, *assumeData.Credentials.Expiration)
	}
	return nil
}

func persistSharedCredentials(c *sts.Credentials, profile string) {
	cred.SaveCredentials(profile, *c.AccessKeyId, *c.SecretAccessKey, *c.SessionToken)
}
