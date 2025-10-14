package currency

import (
	"encoding/xml"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

func LoadCurrencies(filePath string) ValCurs {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// Конвертируем из Windows-1251 в UTF-8
	decoder := charmap.Windows1251.NewDecoder()
	decodedData, err := decoder.Bytes(data)
	if err != nil {
		panic("Ошибка конвертации кодировки: " + err.Error())
	}

	// Просто удаляем XML declaration - самое простое решение
	xmlStr := string(decodedData)
	if idx := strings.Index(xmlStr, "?>"); idx != -1 {
		xmlStr = xmlStr[idx+2:] // берем все после ?>
	}

	var valCurs ValCurs
	err = xml.Unmarshal([]byte(xmlStr), &valCurs)
	if err != nil {
		panic("Error: parsing XML: " + err.Error())
	}

	return valCurs
}
