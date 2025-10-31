package encoder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/identicalaffiliation/task-3/internal/currency"
)

func SaveToJson(path string, data []currency.Currency) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("fail to create directory")
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("fail to create output file")
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("fail to encode JSON")
	}

	return nil
}
