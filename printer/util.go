package printer

import (
	"fmt"
	"strings"

	"github.com/sam701/awstools/colors"
)

func PrintProperties(indent int, pairs ...string) {
	maxKeyLen := 0
	for i := 0; i < len(pairs); i += 2 {
		if len(pairs[i]) > maxKeyLen {
			maxKeyLen = len(pairs[i])
		}
	}

	for i := 0; i < len(pairs); i += 2 {
		prop := pairs[i]
		value := pairs[i+1]

		fmt.Print(strings.Repeat(" ", indent))
		pl := len(prop)
		fmt.Print(colors.Property(prop + ":"))
		fmt.Print(strings.Repeat(" ", maxKeyLen+3-pl))
		fmt.Println(value)

	}
}
