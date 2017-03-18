package cf

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/sam701/tcolor"

	"github.com/sam701/awstools/printer"
	"github.com/urfave/cli"
)

type stackPrintOptions struct {
	printParameters bool
	printOutputs    bool
	printResources  bool
	printTags       bool
}

func createPrintOptions(ctx *cli.Context) *stackPrintOptions {
	return &stackPrintOptions{
		printParameters: ctx.Bool("print-parameters"),
		printOutputs:    ctx.Bool("print-outputs"),
		printTags:       ctx.Bool("print-tags"),
		printResources:  ctx.Bool("print-resources"),
	}
}

var (
	colorStack   = tcolor.New().Foreground(tcolor.BrightWhite).Bold()
	colorSection = tcolor.New().Foreground(tcolor.Yellow).Underline()
)

func printStacks(stacks []*cloudformation.Stack, options *stackPrintOptions) {
	for _, stack := range stacks {
		fmt.Println(tcolor.Colorize(*stack.StackName, colorStack))

		if options.printTags {
			fmt.Println(" ", tcolor.Colorize("Tags", colorSection))
			pairs := []string{}
			for _, tag := range stack.Tags {
				pairs = append(pairs, *tag.Key, *tag.Value)
			}
			printer.PrintProperties(4, pairs...)
		}

		if options.printParameters {
			fmt.Println(" ", tcolor.Colorize("Parameters", colorSection))
			pairs := []string{}
			for _, par := range stack.Parameters {
				pairs = append(pairs, *par.ParameterKey, *par.ParameterValue)
			}
			printer.PrintProperties(4, pairs...)
		}

		if options.printOutputs {
			fmt.Println(" ", tcolor.Colorize("Outputs", colorSection))
			pairs := []string{}
			for _, output := range stack.Outputs {
				pairs = append(pairs, *output.OutputKey, *output.OutputValue)
			}
			printer.PrintProperties(4, pairs...)
		}
	}
}
