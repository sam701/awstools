package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aybabtme/rgbterm"
	"github.com/codegangsta/cli"
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
	client := ec2.New(currentEnvVarSession())
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

	fmt.Println(rgbterm.FgString("Instances:", 255, 255, 255))
	for _, r := range res.Reservations {
		for _, in := range r.Instances {

			var iid string
			if *in.State.Name == "running" {
				iid = rgbterm.FgString(*in.InstanceId, 130, 255, 130)
			} else {
				iid = rgbterm.FgString(*in.InstanceId, 255, 80, 80)
			}
			fmt.Println(
				iid,
				in.LaunchTime.Format("2006-01-02 15:04"),
				*in.InstanceType,
				rgbterm.FgString(flattenString(in.PrivateIpAddress), 80, 80, 255),
			)
			printTags(in, searchPattern)
			fmt.Println()
		}
	}
}
func printTags(instance *ec2.Instance, searchPattern string) {
	for _, tag := range sortedEc2Tags(instance) {
		val := strings.Replace(*tag.Value, searchPattern, rgbterm.FgString(searchPattern, 190, 80, 80), -1)
		fmt.Printf("  %-30s %s\n", *tag.Key+":", val)
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
