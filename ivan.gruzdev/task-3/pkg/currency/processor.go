package currency

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

const (
	filePerm = 0o600
	dirPerm  = 0o755
)

func SortValues(currencies *ValCurs) {
	sort.Slice(currencies.Currencies, func(i, j int) bool {
		return ParseValue(currencies.Currencies[i].Value) > ParseValue(currencies.Currencies[j].Value)
	})
}

func SaveToJSON(filePath string, currencies ValCurs) {
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		panic("Ошибка создания директории: " + err.Error())
	}

	result := make([]OutputCurrency, 0, len(currencies.Currencies))

	for _, currency := range currencies.Currencies {
		object := OutputCurrency{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    ParseValue(currency.Value),
		}
		result = append(result, object)
	}

	data, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		panic(err.Error())
	}

	err = os.WriteFile(filePath, data, filePerm)
	if err != nil {
		panic(err.Error)
	}
}
