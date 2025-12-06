//go:build dev

package config

import (
	"log"

	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devConfig []byte

func init() {
	err := yaml.Unmarshal(devConfig, &cfg)
	if err != nil {
		log.Fatalf("failed to parse dev.yaml: %v", err)
	}
}
