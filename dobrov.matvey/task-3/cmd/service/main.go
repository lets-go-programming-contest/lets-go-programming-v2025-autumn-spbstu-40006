package main

func main() {
	configPath := getConfigPath()

	var (
		cfg  Config
		curs ValCurs
	)

	err := readDataFromConfig(&cfg, configPath)

	if err != nil {
		return
	}

	err = readDataFileNCanGetCurs(&curs, cfg.InputFile)

	if err != nil {
		return
	}

	rates := fillNSortRates(curs)

	fillOutputFile(rates, cfg)
}
