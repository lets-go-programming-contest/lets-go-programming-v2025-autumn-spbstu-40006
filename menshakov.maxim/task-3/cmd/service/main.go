package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valute  []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int    `xml:"NumCode" json:"num_code"`
	CharCode string `xml:"CharCode" json:"char_code"`
	Nominal  int    `xml:"Nominal" json:"-"`
	ValueRaw string `xml:"Value" json:"-"`
	Value float64 `xml:"-" json:"value"`
}

func parseConfig(path string) Config {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("cannot open config file %q: %w", path, err))
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		panic(fmt.Errorf("cannot decode config yaml: %w", err))
	}
	if strings.TrimSpace(cfg.InputFile) == "" || strings.TrimSpace(cfg.OutputFile) == "" {
		panic(fmt.Errorf("config: input-file and output-file must be set"))
	}
	return cfg
}

func parseCBRXML(path string) []Valute {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("cannot open input xml file %q: %w", path, err))
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Errorf("cannot read xml file: %w", err))
	}

	var vc ValCurs
	if err := xml.Unmarshal(data, &vc); err != nil {
		panic(fmt.Errorf("xml unmarshal error: %w", err))
	}

	for i := range vc.Valute {
		raw := strings.TrimSpace(vc.Valute[i].ValueRaw)
		raw = strings.ReplaceAll(raw, " ", "")
		raw = strings.ReplaceAll(raw, ",", ".")
		if raw == "" {
			panic(fmt.Errorf("empty Value for currency %s", vc.Valute[i].CharCode))
		}
		val, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			panic(fmt.Errorf("cannot parse Value %q for %s: %w", vc.Valute[i].ValueRaw, vc.Valute[i].CharCode, err))
		}
		if vc.Valute[i].Nominal == 0 {
			panic(fmt.Errorf("nominal is zero for currency %s", vc.Valute[i].CharCode))
		}
		vc.Valute[i].Value = val / float64(vc.Valute[i].Nominal)
	}

	return vc.Valute
}

func ensureDirForFile(path string) {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(fmt.Errorf("cannot create output directory %q: %w", dir, err))
	}
}

func main() {
	cfgPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()
	if *cfgPath == "" {
		panic("config flag -config is required")
	}

	cfg := parseConfig(*cfgPath)

	vals := parseCBRXML(cfg.InputFile)

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Value > vals[j].Value
	})

	outJSON, err := json.MarshalIndent(vals, "", "  ")
	if err != nil {
		panic(fmt.Errorf("cannot marshal result to json: %w", err))
	}

	ensureDirForFile(cfg.OutputFile)
	if err := os.WriteFile(cfg.OutputFile, outJSON, 0o644); err != nil {
		panic(fmt.Errorf("cannot write output file %q: %w", cfg.OutputFile, err))
	}

	fmt.Printf("Wrote %d currencies to %s\n", len(vals), cfg.OutputFile)
}
