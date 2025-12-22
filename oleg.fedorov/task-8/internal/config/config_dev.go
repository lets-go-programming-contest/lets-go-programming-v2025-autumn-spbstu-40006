//go:build dev

package config

import (
	_ "embed"
)

//go:embed configs/dev.yaml
var configData []byte

type embedLoader struct{}

func (e *embedLoader) Load() ([]byte, error) {
	return configData, nil
}

func NewLoader() Loader {
	return &embedLoader{}
}
