package assume

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/service/sts"
)

var exportPattern = ""
var scriptOutput io.Writer = os.Stdout

func init() {
	shell := path.Base(os.Getenv("SHELL"))
	if shell == "fish" {
		exportPattern = "set -xg %s \"%s\"\n"
	} else {
		exportPattern = "export %s=\"%s\"\n"
	}

}

func printExport(key, value string) {
	fmt.Fprintf(scriptOutput, exportPattern, key, value)
}

func printShellVariables(profileName string, credendials *sts.Credentials) {
	if exportProfile {
		printExportProfile(profileName)
	} else {
		printExportKeyAndToken(credendials)
	}
}

func printExportProfile(profile string) {
	printExport("AWS_PROFILE", profile)
}

func printExportKeyAndToken(cred *sts.Credentials) {
	printExport("AWS_ACCESS_KEY_ID", *cred.AccessKeyId)
	printExport("AWS_SECRET_ACCESS_KEY", *cred.SecretAccessKey)
	printExport("AWS_SESSION_TOKEN", *cred.SessionToken)
}
