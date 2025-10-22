package processing

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveJSON(path string, currencies []Currency) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	if err := json.NewEncoder(file).Encode(currencies); err != nil {
		panic(err)
	}
}
