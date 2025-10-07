package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CommaFloat float64

var (
	ErrEmptyNumber        = errors.New("empty number")
	ErrMultipleSeparators = errors.New("multiple decimal separators")
	ErrInvalidNumber      = errors.New("invalid number")
)

const (
	marshalBufCap     = 32
	marshalFracDigits = 4
	floatBitSize      = 64
)

func (cf *CommaFloat) UnmarshalText(text []byte) error {
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

	*cf = CommaFloat(v)

	return nil
}

func (cf *CommaFloat) MarshalText() ([]byte, error) {
	buf := make([]byte, 0, marshalBufCap)
	buf = strconv.AppendFloat(buf, float64(*cf), 'f', marshalFracDigits, floatBitSize)

	return buf, nil
}
