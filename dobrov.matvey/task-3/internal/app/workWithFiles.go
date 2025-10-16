package app

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

func ReadDataFromConfig(cfg *Config, configPath string) error {
	info, err := os.Stat(configPath)

	if err != nil {
		return fmt.Errorf("file aint exist")
	}

	if info.Size() == 0 {
		return fmt.Errorf("file empty")
	}

	data, err := os.ReadFile(configPath)

	if err != nil {
		return fmt.Errorf("aint read file")
	}

	err = yaml.Unmarshal(data, cfg)

	if err != nil {
		return fmt.Errorf("error with parse config struct")
	}

	return nil
}

func ReadDataFileNCanGetCurs(curs *ValCurs, inputFile string) error {
	f, err := os.Open(inputFile)

	if err != nil {
		return fmt.Errorf("aint can open data file")
	}

	defer f.Close()

	dec := xml.NewDecoder(f)
	dec.CharsetReader = charset.NewReaderLabel

	err = dec.Decode(curs)

	if err != nil {
		return fmt.Errorf("aint can decode data file")
	}

	return nil
}

func FillOutputFile(rates []Rate, cfg Config) error {
	jsonData, err := json.MarshalIndent(rates, "", " ")

	if err != nil {
		return fmt.Errorf("error with create json")
	}

	os.Mkdir(filepath.Dir(cfg.OutputFile), 0755)
	err = os.WriteFile(cfg.OutputFile, jsonData, 0644)

	if err != nil {
		return fmt.Errorf("error with fill json file")
	}

	return nil
}
