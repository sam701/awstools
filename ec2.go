package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sam701/awstools/colors"
	"github.com/sam701/awstools/sess"
	"github.com/sam701/tcolor"
	"github.com/urfave/cli"
)

func actionDescribeEC2(c *cli.Context) error {
	searchPattern := c.Args().First()
	if searchPattern == "" {
		cli.ShowCommandHelp(c, "ec2")
	} else {
		printInstanceStatus(searchPattern)
	}
	return nil
}

func printInstanceStatus(searchPattern string) {
	client := ec2.New(sess.FromEnvVar())
	res, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag-value"),
				Values: []*string{
					aws.String(fmt.Sprintf("*%s*", searchPattern)),
				},
			},
		},
	})
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	fmt.Println(tcolor.Colorize("Instances",
		tcolor.New().Foreground(tcolor.BrightWhite).Bold().Underline()))
	for _, r := range res.Reservations {
		for _, in := range r.Instances {

			var idColor tcolor.Color
			if *in.State.Name == "running" {
				idColor = tcolor.New().Bold().Foreground(tcolor.BrightGreen)
			} else {
				idColor = tcolor.New().Bold().Foreground(tcolor.BrightRed)
			}
			fmt.Println(
				tcolor.Colorize(*in.InstanceId, idColor),
				in.LaunchTime.Format("2006-01-02 15:04"),
				tcolor.Colorize(*in.InstanceType, tcolor.New().Foreground(tcolor.Yellow).Italic()),
				tcolor.Colorize(flattenString(in.PrivateIpAddress), tcolor.New().Foreground(tcolor.BrightBlue)),
				tcolor.Colorize(flattenString(in.PublicIpAddress), tcolor.New().Foreground(tcolor.Cyan)),
			)
			printTags(in, searchPattern)
			fmt.Println()
		}
	}
}
func printTags(instance *ec2.Instance, searchPattern string) {
	for _, tag := range sortedEc2Tags(instance) {
		val := strings.Replace(*tag.Value, searchPattern, colors.Match(searchPattern), -1)
		fmt.Printf("  %s %s\n",
			colors.Property(fmt.Sprintf("%-30s", *tag.Key+":")),
			val)
	}
}

func sortedEc2Tags(instance *ec2.Instance) []*ec2.Tag {
	mm := make(map[string]*ec2.Tag)
	keys := []string{}
	for _, tag := range instance.Tags {
		key := *tag.Key
		mm[key] = tag
		keys = append(keys, key)
	}
	sort.StringSlice(keys).Sort()

	result := []*ec2.Tag{}
	for _, key := range keys {
		result = append(result, mm[key])
	}
	return result
}

func flattenString(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}
