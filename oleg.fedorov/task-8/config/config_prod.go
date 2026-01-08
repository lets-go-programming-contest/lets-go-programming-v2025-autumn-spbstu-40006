//go:build !dev

package config

import _ "embed"

//go:embed configs/prod.yaml
var ConfigFile []byte
