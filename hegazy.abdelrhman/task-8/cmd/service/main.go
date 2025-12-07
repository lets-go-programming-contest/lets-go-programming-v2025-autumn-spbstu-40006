package main

//nolint:gofumpt
import (
	"fmt"

	"task-8/pkg/config"
)

func main() {
	config := config.New()

	fmt.Print(config.Environment, " ", config.LogLevel)
}