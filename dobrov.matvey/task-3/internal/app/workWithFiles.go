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

func ReadDataFromConfig(cfg *Config, configPath string) error {
	info, err := os.Stat(configPath)
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		return err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}

func ReadDataFileNCanGetCurs(curs *ValCurs, inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	dec := xml.NewDecoder(file)
	dec.CharsetReader = charset.NewReaderLabel

	err = dec.Decode(curs)
	if err != nil {
		return err
	}

	return nil
}

func FillOutputFile(rates []Rate, cfg Config) error {
	jsonData, err := json.MarshalIndent(rates, "", " ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(cfg.OutputFile), dirPerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(cfg.OutputFile, jsonData, filePerm)
	if err != nil {
		return err
	}

	return nil
}
