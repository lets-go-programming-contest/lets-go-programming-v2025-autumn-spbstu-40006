package processor

import "fmt"

type CurrencyProcessor struct {
	config *Config
}

func NewCurrencyProcessor() *CurrencyProcessor {
	return &CurrencyProcessor{
		config: nil,
	}
}

func (cp *CurrencyProcessor) Run() error {
	if err := cp.LoadConfig(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	currencies, err := cp.ParseXMLFile(cp.config.InputFile)
	if err != nil {
		return fmt.Errorf("XML parsing error: %w", err)
	}

	if err := cp.SaveToJSON(currencies, cp.config.OutputFile); err != nil {
		return fmt.Errorf("JSON saving error: %w", err)
	}

	return nil
}
