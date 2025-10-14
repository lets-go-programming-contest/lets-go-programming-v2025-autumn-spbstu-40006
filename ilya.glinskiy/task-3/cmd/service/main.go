package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type Input struct {
}

func main() {
	var config Config

	if len(os.Args) < 3 || os.Args[1] != "-config" {
		fmt.Println("Config file must be passed in parameters")

		return
	}

	var configPath string = os.Args[2]

	content, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Couldn't open config file")

		return
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		fmt.Println("Something wrong I can feel it")

		return
	}

	content, err = os.ReadFile(configPath)
}
