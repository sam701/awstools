package sess

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sam701/awstools/config"
)

func FromEnvVar() *session.Session {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = config.Current.DefaultRegion
	}
	return session.New(&aws.Config{
		Region: aws.String(region),
	})
}

func New(name string) *session.Session {
	return session.New(&aws.Config{
		Region:      aws.String(config.Current.DefaultRegion),
		Credentials: credentials.NewSharedCredentials("", name),
	})
}
