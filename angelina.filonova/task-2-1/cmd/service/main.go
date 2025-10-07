package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/filon6/task-2-1/pkg/department"
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
		return 0, fmt.Errorf("Ошибка: не удалось прочитать число")
	}

	text := strings.TrimSpace(scanner.Text())
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("Ошибка: неверный формат числа")
	}
	return value, nil
}

func processDepartment(scanner *bufio.Scanner) error {
	countEmployees, err := readInt(scanner)
	if err != nil {
		return fmt.Errorf("Ошибка: некорректное количество сотрудников")
	}

	dept := department.NewDepartment(countEmployees)

	for range countEmployees {
		if !scanner.Scan() {
			return fmt.Errorf("Ошибка: не удалось прочитать строку с требованием")
		}

		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) != partsCount {
			return fmt.Errorf("Ошибка: неверный формат строки")
		}

		op := parts[0]
		temp, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("Ошибка: неверная температура")
		}

		result := dept.ProcessWorkerRequirement(op, temp)
		fmt.Println(result)
	}

	return nil
}
