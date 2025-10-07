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
	ErrReadNumberFailed     = errors.New("Ошибка: не удалось прочитать число")
	ErrInvalidNumberFormat  = errors.New("Ошибка: неверный формат числа")
	ErrInvalidEmployeeCount = errors.New("Ошибка: некорректное количество сотрудников")
	ErrReadRequirement      = errors.New("Ошибка: не удалось прочитать строку с требованием")
	ErrInvalidLineFormat    = errors.New("Ошибка: неверный формат строки")
	ErrInvalidTemperature   = errors.New("Ошибка: неверная температура")
)

const partsCount = 2

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	countDepartments, err := readInt(scanner)
	if err != nil {
		fmt.Println("Ошибка: некорректное количество отделов")
		return
	}

	for range countDepartments {
		if err := processDepartment(scanner); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func readInt(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, ErrReadNumberFailed
	}

	text := strings.TrimSpace(scanner.Text())
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, ErrInvalidNumberFormat
	}

	return value, nil
}

func processDepartment(scanner *bufio.Scanner) error {
	countEmployees, err := readInt(scanner)
	if err != nil {
		return ErrInvalidEmployeeCount
	}

	dept := department.NewDepartment(countEmployees)

	for range countEmployees {
		if !scanner.Scan() {
			return ErrReadRequirement
		}

		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) != partsCount {
			return ErrInvalidLineFormat
		}

		operand := parts[0]
		temp, err := strconv.Atoi(parts[1])
		if err != nil {
			return ErrInvalidTemperature
		}

		result := dept.ProcessWorkerRequirement(operand, temp)
		fmt.Println(result)
	}

	return nil
}
