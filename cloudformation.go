package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aybabtme/rgbterm"
	"github.com/codegangsta/cli"
)

func printStacks(c *cli.Context) error {
	pattern := c.String("search")
	if pattern == "" {
		cli.ShowCommandHelp(c, "cloudformation")
	} else {
		cf := cloudformation.New(currentEnvVarSession())
		out, err := cf.DescribeStacks(&cloudformation.DescribeStacksInput{})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		for _, stack := range out.Stacks {
			if strings.Contains(*stack.StackName, pattern) {
				formatStack(stack)
			}
		}
	}
	return nil
}

func formatStack(stack *cloudformation.Stack) {
	fmt.Printf("%s %s\n",
		rgbterm.FgString(fmt.Sprintf("%-30s", *stack.StackName), 255, 255, 255),
		rgbterm.FgString(*stack.StackStatus, 130, 255, 130))
	fmt.Println("  Parameters:")
	for _, par := range stack.Parameters {
		fmt.Printf("    %-35s %s\n", *par.ParameterKey, *par.ParameterValue)
	}
	fmt.Println("  Outputs:")
	for _, out := range stack.Outputs {
		fmt.Printf("    %-35s %s\n", *out.OutputKey, *out.OutputValue)
	}
	fmt.Println()
}
