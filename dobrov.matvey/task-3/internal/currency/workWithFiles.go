package currency

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

func ReadDataFileNCanGetCurs(curs *ValCurs, inputFile string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("read %q: %w", inputFile, err)
	}

	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.CharsetReader = charset.NewReaderLabel

	if err := dec.Decode(curs); err != nil {
		return fmt.Errorf("xml decode %q: %w", inputFile, err)
	}

	return nil
}

func FillOutputFile(currency []Currency, outputPath string) error {
	jsonData, err := json.MarshalIndent(currency, "", " ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("mkdir %q: %w", dir, err)
	}

	if err := os.WriteFile(outputPath, jsonData, filePerm); err != nil {
		return fmt.Errorf("write %q: %w", outputPath, err)
	}

	return nil
}
