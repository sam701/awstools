package assume

import (
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:      "assume",
	Usage:     "assume role on a specified account",
	ArgsUsage: "<account name> <role to assume>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "export, e",
			Usage: "export shell script to `PATH`",
		},
		cli.BoolFlag{
			Name:  "export-profile, p",
			Usage: "export AWS_PROFILE instead of token variables (this will become default behavior later)",
		},
		cli.BoolFlag{
			Name:  "export-token",
			Usage: "export AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_SESSION_TOKEN variables (this is currently default behavior)",
		},
		cli.IntFlag{
			Name:  "reuse-credentials",
			Usage: "reuse current credentials if they are valid for at least `MINUTES`",
		},
	},
	Action: assumeRoleAction,
}
