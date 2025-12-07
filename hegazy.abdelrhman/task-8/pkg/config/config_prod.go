//go:build !dev

package config

import (
	_ "embed"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodConfig []byte

func init() {
	err := yaml.Unmarshal(prodConfig, &cfg)
	if err != nil {
		log.Fatalf("failed to parse prod.yaml: %v", err)
	}
}