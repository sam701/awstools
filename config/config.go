package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/naoina/toml"
	"github.com/sam701/tcolor"
)

var Current *Configuration

type Configuration struct {
	DefaultRegion              string
	DefaultKmsKey              string
	AutoRotateMainAccountKey   bool // Deprecated: use KeyRotationIntervalMinutes
	KeyRotationIntervalMinutes int

	ReuseCredentialsIfValidForMinutes int

	Profiles struct {
		MainAccount           string
		MainAccountMfaSession string
	}
	Accounts map[string]string
}

func Read(filePath string) {
	if filePath == "" {
		filePath = path.Join(os.Getenv("HOME"), ".config", "awstools", "config.toml")
	}
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()

	var c Configuration
	c.KeyRotationIntervalMinutes = 60 * 24 * 7 // 1 week
	c.ReuseCredentialsIfValidForMinutes = 120  // turned off

	err = toml.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	Current = &c

	checkDeprecatedValues()
}

func checkDeprecatedValues() {
	if Current.AutoRotateMainAccountKey {
		fmt.Println()
		fmt.Println(
			tcolor.Colorize("WARNING:", tcolor.New().Foreground(tcolor.BrightRed).Bold()),
			tcolor.Colorize("autoRotateMainAccountKey", tcolor.New().Foreground(tcolor.Red).Underline()),
			tcolor.Colorize("in your config.toml is deprecated, use keyRotationIntervalMinutes instead", tcolor.New().Foreground(tcolor.Red)),
		)
		fmt.Println()
	}
}
