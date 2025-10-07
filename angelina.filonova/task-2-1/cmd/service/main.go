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

const partsCount = 2

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	countDepartments, err := readInt(scanner)
	if err != nil {
		fmt.Println("Ошибка: некорректное количество отделов")

		return
	}

	for i := 0; i < countDepartments; i++ {
		if err := processDepartment(scanner); err != nil {
			fmt.Println(err)

			return
		}
	}
}

func readInt(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, errors.New("Ошибка: не удалось прочитать число")
	}

	text := strings.TrimSpace(scanner.Text())
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, errors.New("Ошибка: неверный формат числа")
	}

	return value, nil
}

func processDepartment(scanner *bufio.Scanner) error {
	countEmployees, err := readInt(scanner)
	if err != nil {
		return errors.New("Ошибка: некорректное количество сотрудников")
	}

	dept := department.NewDepartment(countEmployees)

	for i := 0; i < countEmployees; i++ {
		if !scanner.Scan() {
			return errors.New("Ошибка: не удалось прочитать строку с требованием")
		}

		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)

		if len(parts) != partsCount {
			return errors.New("Ошибка: неверный формат строки")
		}

		operand := parts[0]
		temp, err := strconv.Atoi(parts[1])
		if err != nil {
			return errors.New("Ошибка: неверная температура")
		}

		result := dept.ProcessWorkerRequirement(operand, temp)
		fmt.Println(result)
	}

	return nil
}
