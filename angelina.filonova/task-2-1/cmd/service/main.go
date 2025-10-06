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

	countDepartments, err := readInt(scanner, "количество отделов")
	if err != nil {
		fmt.Println(err)
		return
	}

	for range countDepartments {
		if err := processDepartment(scanner); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func readInt(scanner *bufio.Scanner, name string) (int, error) {
	if !scanner.Scan() {
		return 0, fmt.Errorf("не удалось прочитать %s", name)
	}

	valueStr := strings.TrimSpace(scanner.Text())
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("неверный формат числа %s: %v", name, err)
	}
	return value, nil
}

func processDepartment(scanner *bufio.Scanner) error {
	countEmployees, err := readInt(scanner, "количество сотрудников")
	if err != nil {
		return err
	}

	dept := department.NewDepartment(countEmployees)

	for range countEmployees {
		if !scanner.Scan() {
			return fmt.Errorf("не удалось прочитать строку с требованием")
		}

		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) != partsCount {
			return fmt.Errorf("неверный формат строки: %q", line)
		}

		op := parts[0]
		temp, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("неверная температура: %q (%v)", parts[1], err)
		}

		result := dept.ProcessWorkerRequirement(op, temp)
		fmt.Println(result)
	}

	return nil
}
