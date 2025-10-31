package currency

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func LoadCurrencies(filePath string) (ValCurs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ValCurs{}, fmt.Errorf("failed to read file: %w", err)
	}

	decoder := charmap.Windows1251.NewDecoder()

	decodedData, err := decoder.Bytes(data)
	if err != nil {
		return ValCurs{}, fmt.Errorf("encoding conversion error: %w", err)
	}

	xmlStr := string(decodedData)
	if idx := strings.Index(xmlStr, "?>"); idx != -1 {
		xmlStr = xmlStr[idx+2:]
	}

	var valCurs ValCurs

	err = xml.Unmarshal([]byte(xmlStr), &valCurs)
	if err != nil {
		return ValCurs{}, fmt.Errorf("failed to parse XML: %w", err)
	}

	return valCurs, nil
}
