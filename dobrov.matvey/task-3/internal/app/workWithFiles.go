package app

import (
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

const (
	dirPerm  = 0o755
	filePerm = 0o600
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
		return errReadConfig
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return errReadConfig
	}

	return nil
}

func ReadDataFileNCanGetCurs(curs *ValCurs, inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return errOpenData
	}

	defer func() { _ = file.Close() }()

	dec := xml.NewDecoder(file)
	dec.CharsetReader = charset.NewReaderLabel

	err = dec.Decode(curs)
	if err != nil {
		return errDecodeData
	}

	return nil
}

func FillOutputFile(rates []Rate, cfg Config) error {
	jsonData, err := json.MarshalIndent(rates, "", " ")
	if err != nil {
		return errCreateJSON
	}

	err = os.MkdirAll(filepath.Dir(cfg.OutputFile), dirPerm)
	if err != nil {
		return fmt.Errorf("mkdir error %w", err)
	}

	err = os.WriteFile(cfg.OutputFile, jsonData, filePerm)
	if err != nil {
		return errWriteOutputJSON
	}

	return nil
}
