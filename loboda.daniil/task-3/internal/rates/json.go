package rates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

func SaveJSON(path string, list []Currency) error {
	out := make([]Currency, len(list))
	copy(out, list)

	sort.Slice(out, func(i, j int) bool {
		return out[i].ValueNum > out[j].ValueNum
	})

	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}

	data, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(path, data, filePerm); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return nil
}
