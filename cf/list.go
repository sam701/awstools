package cf

import (
	"log"
	"sort"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/sam701/awstools/sess"
	"github.com/urfave/cli"
)

func listStacks(ctx *cli.Context) error {
	cfClient = cloudformation.New(sess.FromEnvVar())

	filter := createStackFilter(
		ctx.String("name"),
		ctx.String("name-substring"),
		ctx.StringSlice("tag"),
	)

	printOptions := createPrintOptions(ctx)

	stacks := getStacks(filter)
	printStacks(stacks, printOptions)
	return nil
}

func getStacks(filter *stackFilter) []*cloudformation.Stack {
	stacks := []*cloudformation.Stack{}

	iterateStacks(func(stack *cloudformation.Stack) bool {
		if filter.match(stack) {
			stacks = append(stacks, stack)
		}
		return true
	})

	sort.Slice(stacks, func(a, b int) bool {
		return *stacks[a].StackName < *stacks[b].StackName
	})
	return stacks
}

func iterateStacks(handler func(*cloudformation.Stack) bool) {
	var nextToken *string
	for {
		out, err := cfClient.DescribeStacks(&cloudformation.DescribeStacksInput{
			NextToken: nextToken,
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		nextToken = out.NextToken

		for _, stack := range out.Stacks {
			proceed := handler(stack)
			if !proceed {
				return
			}
		}

		if nextToken == nil {
			break
		}
	}
}
