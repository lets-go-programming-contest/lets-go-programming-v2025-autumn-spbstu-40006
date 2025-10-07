package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type CommaFloat float64

func (cf *CommaFloat) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))
	if s == "" {
		return fmt.Errorf("empty number")
	}

	s = strings.Replace(s, ",", ".", 1)
	if strings.Count(s, ",") > 0 {
		return fmt.Errorf("multiple decimal separators in %q", string(text))
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q: %w", string(text), err)
	}

	*cf = CommaFloat(v)
	return nil
}
