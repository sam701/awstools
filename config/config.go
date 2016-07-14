package config

import (
	"log"
	"os"
	"path"

	"github.com/naoina/toml"
)

var Current *Configuration

type Configuration struct {
	DefaultRegion            string
	DefaultKmsKey            string
	AutoRotateMainAccountKey bool

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
	err = toml.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	Current = &c
}
