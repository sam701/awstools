package cf

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/sam701/tcolor"
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
		tcolor.Colorize(fmt.Sprintf("%-30s", *stack.StackName), tcolor.New().Foreground(tcolor.BrightWhite).Bold()),
		tcolor.Colorize(*stack.StackStatus, tcolor.New().Foreground(tcolor.BrightGreen)))
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
