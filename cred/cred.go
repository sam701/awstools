package cred

import (
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/service/sts"
)

func credentialsFilePath() string {
	return path.Join(os.Getenv("HOME"), ".aws/credentials")
}

func configFilePath() string {
	return path.Join(os.Getenv("HOME"), ".aws/config")
}

func SaveCredentials(profile, keyId, secret, token string) {
	updatePropsGroup(credentialsFilePath(), newCredentialsGroup(profile, keyId, secret, token))
}

func SetProfileRegion(profile, region string) {
	updatePropsGroup(configFilePath(), &propsGroup{
		name: "profile '" + profile + "'",
		lines: []*propertyLine{
			&propertyLine{"region", region},
		},
	})
}

func newCredentialsGroup(profile, keyId, secret, token string) *propsGroup {
	g := &propsGroup{
		name: profile,
		lines: []*propertyLine{
			&propertyLine{"aws_access_key_id", keyId},
			&propertyLine{"aws_secret_access_key", secret},
		},
	}
	if token != "" {
		g.lines = append(g.lines, &propertyLine{"aws_session_token", token})
	}
	return g
}

func GetMainAccountKeyId(profileName string) string {
	cf := readCredentials()
	for _, g := range cf.groups {
		if g.name == profileName {
			for _, l := range g.lines {
				if l.key == "aws_access_key_id" {
					return l.value
				}
			}
		}
	}
	log.Fatalf("Credentials file does not contain profile '%s'\n", profileName)
	return ""
}

func readCredentials() *propsFile {
	f, err := os.Open(credentialsFilePath())
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	return readPropsFile(f)
}

func GetCredentials(profileName string) *sts.Credentials {
	cf := readCredentials()
	for _, g := range cf.groups {
		if g.name == profileName {
			cr := &sts.Credentials{}
			for _, n := range g.lines {
				switch n.key {
				case "aws_access_key_id":
					cr.SetAccessKeyId(n.value)
				case "aws_secret_access_key":
					cr.SetSecretAccessKey(n.value)
				case "aws_session_token":
					cr.SetSessionToken(n.value)
				}
			}
			return cr
		}
	}
	return nil
}
