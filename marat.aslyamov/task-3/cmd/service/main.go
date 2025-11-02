package main

import (
	"github.com/tuesdayy1/task-3/internal/processor"
)

func main() {
	app := processor.NewCurrencyProcessor()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
