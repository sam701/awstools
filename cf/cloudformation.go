package cf

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

var cfClient *cloudformation.CloudFormation

func awsToString(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}
