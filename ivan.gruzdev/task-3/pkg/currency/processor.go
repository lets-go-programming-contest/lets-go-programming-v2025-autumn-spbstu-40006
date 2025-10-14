package currency

import (
	"encoding/json"
	"os"
	"sort"
)

func SortValues(currencies *ValCurs) {
	sort.Slice(currencies.Currencies, func(i, j int) bool {
		return ParseValue(currencies.Currencies[i].Value) > ParseValue(currencies.Currencies[j].Value)
	})
}

func SaveToJson(filePath string, currencies ValCurs) {
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
		panic(err)
	}

	os.WriteFile(filePath, data, 0644)

}
