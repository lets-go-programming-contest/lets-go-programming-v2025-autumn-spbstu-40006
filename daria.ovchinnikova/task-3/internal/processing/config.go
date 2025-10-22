package processing

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cnfg Config
	if err := yaml.Unmarshal(data, &cnfg); err != nil {
		panic(err)
	}

	return cnfg
}
