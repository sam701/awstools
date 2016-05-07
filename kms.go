package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/codegangsta/cli"
)

func kmsAction(c *cli.Context) error {
	txt := c.Args().First()
	if txt == "" {
		cli.ShowCommandHelp(c, "kms")
		return nil
	}
	if c.Bool("decrypt") {
		cl := kms.New(currentEnvVarSession())
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
		fmt.Println("Decrypted:", string(out.Plaintext))
	} else {
		cli.ShowCommandHelp(c, "kms")
	}
	return nil
}
