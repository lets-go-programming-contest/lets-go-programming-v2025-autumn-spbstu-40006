package rates

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Decimal float64

var (
	ErrEmptyNumber        = errors.New("empty number")
	ErrMultipleSeparators = errors.New("multiple decimal separators")
	ErrInvalidNumber      = errors.New("invalid number")
)

func (cf *Decimal) UnmarshalText(text []byte) error {
	str := strings.TrimSpace(string(text))
	if str == "" {
		return ErrEmptyNumber
	}

	str = strings.Replace(str, ",", ".", 1)
	if strings.Count(str, ",") > 0 {
		return fmt.Errorf("%w: %q", ErrMultipleSeparators, text)
	}

	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("%w: %q: %w", ErrInvalidNumber, text, err)
	}

	*cf = Decimal(v)

	return nil
}
