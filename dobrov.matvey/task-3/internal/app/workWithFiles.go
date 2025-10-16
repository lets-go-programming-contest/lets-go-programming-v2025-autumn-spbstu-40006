package app

import (
	"errors"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
)

var (
	errNoFile          = errors.New("input file does not exist")
	errEmptyFile       = errors.New("input file is empty")
	errReadConfig      = errors.New("read config failed")
	errOpenData        = errors.New("open data file failed")
	errDecodeData      = errors.New("decode data file failed")
	errCreateJSON      = errors.New("create json failed")
	errWriteOutputJSON = errors.New("write json file failed")
)

func ReadDataFromConfig(cfg *Config, configPath string) error {
	info, err := os.Stat(configPath)

	if err != nil {
		return errNoFile
	}

	if info.Size() == 0 {
		return errEmptyFile
	}

	data, err := os.ReadFile(configPath)

	if err != nil {
		return fmt.Errorf("%w: %v", errReadConfig, err)
	}

	err = yaml.Unmarshal(data, cfg)

	if err != nil {
		return fmt.Errorf("%w: %v", errReadConfig, err)
	}

	return nil
}

func ReadDataFileNCanGetCurs(curs *ValCurs, inputFile string) error {
	file, err := os.Open(inputFile)

	if err != nil {
		return fmt.Errorf("%w: %v", errOpenData, err)
	}

	defer func() { _ = file.Close() }()

	dec := xml.NewDecoder(file)
	dec.CharsetReader = charset.NewReaderLabel

	err = dec.Decode(curs)

	if err != nil {
		return fmt.Errorf("%w: %v", errDecodeData, err)
	}

	return nil
}

func FillOutputFile(rates []Rate, cfg Config) error {
	jsonData, err := json.MarshalIndent(rates, "", " ")

	if err != nil {
		return fmt.Errorf("%w: %v", errCreateJSON, err)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.OutputFile), 0o755); err != nil {
		return err
	}
	err = os.WriteFile(cfg.OutputFile, jsonData, 0o600)

	if err != nil {
		return fmt.Errorf("%w: %v", errWriteOutputJSON, err)
	}

	return nil
}
