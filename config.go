package main

import (
	"log"
	"os"
	"path"

	"github.com/naoina/toml"
)

type configuration struct {
	Profiles struct {
		Bastion    string
		BastionMfa string
	}
	DefaultRegion string
	Accounts      map[string]string
}

func readConfig(filePath string) *configuration {
	if filePath == "" {
		filePath = path.Join(os.Getenv("HOME"), ".config", "awstools", "config.toml")
	}
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("ERROR", err)
	}
	defer f.Close()

	var c configuration
	err = toml.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return &c
}
