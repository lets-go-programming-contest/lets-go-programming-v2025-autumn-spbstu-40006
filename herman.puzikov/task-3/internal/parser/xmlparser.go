package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Segfault-chan/task-3/internal/rates"
	"golang.org/x/net/html/charset"
)

func ReadXML(filepath string) (*rates.ExchangeRate, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel // for windows-1251 encoding of the xml

	var exchRates rates.ExchangeRate
	if err := decoder.Decode(&exchRates); err != nil {
		return nil, fmt.Errorf("error decoding XML: %w", err)
	}

	return &exchRates, nil
}

func WriteXML(list []rates.Currency, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("couldn't create a directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		return fmt.Errorf("couldn't open/create a file: %w", err)
	}

	encoder := xml.NewEncoder(file)

	if err := encoder.Encode(list); err != nil {
		return fmt.Errorf("problem while writing json: %w", err)
	}

	return nil
}
