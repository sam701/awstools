package cf

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/urfave/cli"

	"github.com/sam701/awstools/printer"
	"github.com/sam701/tcolor"
)

type stackPrintOptions struct {
	printParameters bool
	printOutputs    bool
	printResources  bool
	printTags       bool

	resourceTypes []string
}

func createPrintOptions(ctx *cli.Context) *stackPrintOptions {
	return &stackPrintOptions{
		printParameters: ctx.Bool("print-parameters"),
		printOutputs:    ctx.Bool("print-outputs"),
		printTags:       ctx.Bool("print-tags"),
		printResources:  ctx.Bool("print-resources"),

		resourceTypes: ctx.StringSlice("resource-type"),
	}
}

var (
	colorStack        = tcolor.New().Foreground(tcolor.BrightWhite).Bold()
	colorStackRef     = tcolor.New().Foreground(tcolor.Cyan).Italic()
	colorSection      = tcolor.New().Foreground(tcolor.Yellow).Underline()
	colorResourceType = tcolor.New().Foreground(tcolor.Magenta)
)

func printStacks(stacks []*cloudformation.Stack, options *stackPrintOptions) {
	if len(options.resourceTypes) > 0 {
		printResourceTypes(stacks, options.resourceTypes)
		return
	}

	for i, stack := range stacks {
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

		if options.printResources {
			fmt.Println(" ", tcolor.Colorize("Resources", colorSection))
			res := getStackResources(*stack.StackName, []string{})
			lastType := ""
			for _, r := range res {
				if lastType != *r.ResourceType {
					fmt.Println("   ", tcolor.Colorize(*r.ResourceType, colorResourceType))
				}
				printStackResource(6, r)
				lastType = *r.ResourceType
			}
		}

		if i < len(stacks)-1 {
			fmt.Println("")
		}
	}
}

func printStackResource(indent int, resource *cloudformation.StackResourceSummary) {
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Println(*resource.PhysicalResourceId)
}

func printResourceTypes(stacks []*cloudformation.Stack, resourceTypes []string) {
	type resource struct {
		summary   *cloudformation.StackResourceSummary
		stackName string
	}
	res := []*resource{}

	for _, stack := range stacks {
		r := getStackResources(*stack.StackName, resourceTypes)
		for _, rs := range r {
			res = append(res, &resource{rs, *stack.StackName})
		}
	}

	sort.Slice(res, func(a, b int) bool {
		ra := *res[a].summary.ResourceType
		rb := *res[b].summary.ResourceType
		if ra == rb {
			return *res[a].summary.PhysicalResourceId < *res[b].summary.PhysicalResourceId
		}
		return ra < rb
	})

	lastType := ""
	for _, r := range res {
		if lastType != *r.summary.ResourceType {
			fmt.Println(tcolor.Colorize(*r.summary.ResourceType, colorResourceType))
		}
		fmt.Printf("  %s (%s)\n", *r.summary.PhysicalResourceId,
			tcolor.Colorize(r.stackName, colorStackRef))

		lastType = *r.summary.ResourceType
	}

}
