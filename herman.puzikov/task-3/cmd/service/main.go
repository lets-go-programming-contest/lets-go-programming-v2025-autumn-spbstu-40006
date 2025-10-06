package main

import (
	"fmt"
	"os"

	"github.com/Segfault-chan/task-3/internal/utils"
)

func main() {
	var files = make(map[string]string, 2)

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Invalid number of args: %v\nThe correct usage is: -config <path-to-config>\n", os.Args)

		return
	}

	if os.Args[1] != "-config" {
		fmt.Fprint(os.Stderr, "Invalid first operand.\nThe correct usage is: -config <path-to-config>\n", os.Args)

		return
	}

	if err := utils.ParseYAML(&files, os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error while parsing the yaml file: %v", err)

		return
	}
}
