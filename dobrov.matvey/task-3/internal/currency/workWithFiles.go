package currency

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"

	"golang.org/x/net/html/charset"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

func ReadDataFileNCanGetCurs(curs *ValCurs, inputFile string) (err error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	defer file.Close()

	dec := xml.NewDecoder(file)
	dec.CharsetReader = charset.NewReaderLabel

	err = dec.Decode(curs)
	if err != nil {
		return err
	}

	return dec.Decode(curs)
}

func FillOutputFile(currency []Currency, outputPath string) error {
	jsonData, err := json.MarshalIndent(currency, "", " ")
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(outputPath), dirPerm); err != nil {
		return err
	}
	return os.WriteFile(outputPath, jsonData, filePerm)
}
