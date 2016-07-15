package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aybabtme/rgbterm"
	"github.com/naoina/toml"
)

var Current *Configuration

type Configuration struct {
	DefaultRegion              string
	DefaultKmsKey              string
	AutoRotateMainAccountKey   bool // Deprecated: use KeyRotationIntervalMinutes
	KeyRotationIntervalMinutes int

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
		fmt.Println(rgbterm.FgString("autoRotateMainAccountKey in your config.toml is deprecated, use keyRotationIntervalMinutes instead", 255, 130, 130))
		fmt.Println()
	}
}
