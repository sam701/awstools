package cf

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aybabtme/rgbterm"
)

func searchStack(pattern string) {
	out, err := cfClient.DescribeStacks(&cloudformation.DescribeStacksInput{})
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	for _, stack := range out.Stacks {
		if strings.Contains(*stack.StackName, pattern) {
			formatStack(stack)
		}
	}
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
