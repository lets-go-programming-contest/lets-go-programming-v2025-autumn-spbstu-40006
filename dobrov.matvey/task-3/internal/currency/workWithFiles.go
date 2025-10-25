package currency

import (
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

func ReadDataFileNCanGetCurs(curs *ValCurs, inputFile string) {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(fmt.Errorf("open %q: %w", inputFile, err))
	}

	defer file.Close()

	dec := xml.NewDecoder(file)
	dec.CharsetReader = charset.NewReaderLabel

	err = dec.Decode(curs)
	if err != nil {
		panic(fmt.Errorf("xml decode %q: %w", inputFile, err))
	}
}

func FillOutputFile(currency []Currency, outputPath string) {
	jsonData, err := json.MarshalIndent(currency, "", " ")
	if err != nil {
		panic(fmt.Errorf("json marshal: %w", err))
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(fmt.Errorf("mkdir %q: %w", dir, err))
	}

	if err := os.WriteFile(outputPath, jsonData, 0o644); err != nil {
		panic(fmt.Errorf("write %q: %w", outputPath, err))
	}
}
