package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func currentEnvVarSession() *session.Session {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = theConfig.DefaultRegion
	}
	return session.New(&aws.Config{
		Region: aws.String(region),
	})
}

func newSession(name string) *session.Session {
	return session.New(&aws.Config{
		Region:      aws.String(theConfig.DefaultRegion),
		Credentials: credentials.NewSharedCredentials("", name),
	})
}
