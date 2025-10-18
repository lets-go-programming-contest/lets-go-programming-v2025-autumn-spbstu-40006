package app

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

func ReadDataFromConfig(cfg *Config, configPath string) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err.Error())
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic(err.Error())
	}
}

func ReadDataFileNCanGetCurs(curs *ValCurs, inputFile string) {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err.Error())
	}

	defer func() { _ = file.Close() }()

	dec := xml.NewDecoder(file)
	dec.CharsetReader = charset.NewReaderLabel

	err = dec.Decode(curs)
	if err != nil {
		panic(err.Error())
	}
}

func FillOutputFile(rates []Rate, cfg Config) {
	jsonData, err := json.MarshalIndent(rates, "", " ")
	if err != nil {
		panic(err.Error())
	}

	err = os.MkdirAll(filepath.Dir(cfg.OutputFile), dirPerm)
	if err != nil {
		panic(err.Error())
	}

	err = os.WriteFile(cfg.OutputFile, jsonData, filePerm)
	if err != nil {
		panic(err.Error())
	}
}
