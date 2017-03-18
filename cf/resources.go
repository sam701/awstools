package cf

import (
	"log"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func getStackResources(stackName string, resourceTypes []string) []*cloudformation.StackResourceSummary {
	var nextToken *string = nil
	res := []*cloudformation.StackResourceSummary{}
	for {
		out, err := cfClient.ListStackResources(&cloudformation.ListStackResourcesInput{
			NextToken: nextToken,
			StackName: aws.String(stackName),
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		nextToken = out.NextToken

		for _, s := range out.StackResourceSummaries {
			if len(resourceTypes) > 0 {
				for _, rt := range resourceTypes {
					if rt == *s.ResourceType {
						res = append(res, s)
						break
					}
				}
			} else {
				res = append(res, s)
			}
		}

		if nextToken == nil {
			break
		}
	}

	sort.Slice(res, func(a, b int) bool {
		return *res[a].ResourceType < *res[b].ResourceType
	})

	return res
}
