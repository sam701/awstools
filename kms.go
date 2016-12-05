package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/sam701/awstools/colors"
	"github.com/sam701/awstools/config"
	"github.com/sam701/awstools/sess"
	"github.com/sam701/tcolor"
	"github.com/urfave/cli"
)

func kmsAction(c *cli.Context) error {
	txt := c.Args().First()

	if c.Bool("list-keys") {
		listKmsKeys()
		return nil
	}

	if txt == "" {
		cli.ShowCommandHelp(c, "kms")
		return nil
	}
	quiet := c.Bool("quiet")
	cl := kms.New(sess.FromEnvVar())
	if c.Bool("decrypt") {
		bb, err := base64.StdEncoding.DecodeString(txt)
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		out, err := cl.Decrypt(&kms.DecryptInput{
			CiphertextBlob: bb,
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		if !quiet {
			fmt.Print("Decrypted: ")
		}
		fmt.Println(string(out.Plaintext))
	} else if c.Bool("encrypt") {
		keyId := c.String("key-id")
		if keyId == "" {
			keyId = config.Current.DefaultKmsKey
		}
		if keyId == "" {
			log.Fatalln("No key-id provided")
		}
		out, err := cl.Encrypt(&kms.EncryptInput{
			KeyId:     aws.String(keyId),
			Plaintext: []byte(txt),
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		if !quiet {
			fmt.Print("Encrypted: ")
		}
		fmt.Println(base64.StdEncoding.EncodeToString(out.CiphertextBlob))
	} else {
		cli.ShowCommandHelp(c, "kms")
	}
	return nil
}

func listKmsKeys() {
	cl := kms.New(sess.FromEnvVar())

	out, err := cl.ListAliases(&kms.ListAliasesInput{})
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	for _, v := range out.Aliases {
		res, err := cl.DescribeKey(&kms.DescribeKeyInput{
			KeyId: v.AliasArn,
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		md := res.KeyMetadata

		fmt.Println(tcolor.Colorize(*md.Arn, tcolor.New().Foreground(tcolor.BrightGreen)))
		fmt.Println(formatProp("Alias"), tcolor.Colorize(*v.AliasName, tcolor.New().Foreground(tcolor.BrightRed)))
		fmt.Println(formatProp("Description"), *md.Description)
		fmt.Println(formatProp("Created"), *md.CreationDate)
		fmt.Println(formatProp("Enabled"), *md.Enabled)
		fmt.Println(formatProp("Usage"), *md.KeyUsage)
		fmt.Println()
	}
}

func formatProp(prop string) string {
	return colors.Property(fmt.Sprintf("  %-15s", prop+":"))
}
