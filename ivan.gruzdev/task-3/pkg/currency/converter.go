package currency

import (
	"strconv"
	"strings"
)

func ParseValue(valueStr string) float64 {
	valueStr = strings.Replace(valueStr, ",", ".", 1)

	valueF, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		panic("Ошибка преобразования числа: " + valueStr)
	}

	return valueF
}
