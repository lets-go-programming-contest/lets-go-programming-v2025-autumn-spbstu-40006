package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/filon6/task-2-1/pkg/department"
)

var (
	ErrReadNumberFailed     = errors.New("failed to read number")
	ErrInvalidNumberFormat  = errors.New("invalid number format")
	ErrInvalidEmployeeCount = errors.New("invalid number of employees")
	ErrReadRequirement      = errors.New("failed to read requirement line")
	ErrInvalidLineFormat    = errors.New("invalid line format")
	ErrInvalidTemperature   = errors.New("invalid temperature value")
)

const partsCount = 2

func main() {
	var (
		scanner          = bufio.NewScanner(os.Stdin)
		countDepartments int
		err              error
	)

	countDepartments, err = readInt(scanner)
	if err != nil {
		fmt.Println("invalid number of departments")

		return
	}

	for range countDepartments {
		err = processDepartment(scanner)
		if err != nil {
			fmt.Println(err)

			return
		}
	}
}

func readInt(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, ErrReadNumberFailed
	}

	var (
		text  = strings.TrimSpace(scanner.Text())
		value int
		err   error
	)

	value, err = strconv.Atoi(text)
	if err != nil {
		return 0, ErrInvalidNumberFormat
	}

	return value, nil
}

func processDepartment(scanner *bufio.Scanner) error {
	var (
		countEmployees int
		err            error
	)

	countEmployees, err = readInt(scanner)
	if err != nil {
		return ErrInvalidEmployeeCount
	}

	dept := department.NewDepartment(countEmployees)

	for range countEmployees {
		if !scanner.Scan() {
			return ErrReadRequirement
		}

		var (
			line  = strings.TrimSpace(scanner.Text())
			parts = strings.Fields(line)
		)

		if len(parts) != partsCount {
			return ErrInvalidLineFormat
		}

		operand := parts[0]

		if temp, err := strconv.Atoi(parts[1]); err != nil {
			return ErrInvalidTemperature
		} else {
			result := dept.ProcessWorkerRequirement(operand, temp)
			fmt.Println(result)
		}
	}

	return nil
}
