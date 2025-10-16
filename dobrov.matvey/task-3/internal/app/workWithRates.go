package main

import (
	"sort"
	"strconv"
	"strings"
)

func fillNSortRates(curs ValCurs) []Rate {
	var rates []Rate

	for _, v := range curs.Valute {
		s := strings.ReplaceAll(v.ValueRaw, ",", ".")
		val, err := strconv.ParseFloat(s, 64)

		if err != nil {
			continue
		}

		if v.Nominal > 0 {
			val = val / float64(v.Nominal)
		}

		rate := Rate{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    val,
		}

		rates = append(rates, rate)
	}

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Value > rates[j].Value
	})

	return rates
}
