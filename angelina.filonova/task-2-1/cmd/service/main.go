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

	if !scanner.Scan() {
		fmt.Println("Не удалось прочитать количество отделов")

		return
	}

	countDepartmentsStr := strings.TrimSpace(scanner.Text())

	countDepartments, err := strconv.Atoi(countDepartmentsStr)
	if err != nil {
		fmt.Println("Неверный формат числа отделов:", err)

		return
	}

	for range countDepartments {
		if !scanner.Scan() {
			fmt.Println("Не удалось прочитать количество сотрудников")

			return
		}
		countEmployeesStr := strings.TrimSpace(scanner.Text())

		countEmployees, err := strconv.Atoi(countEmployeesStr)
		if err != nil {
			fmt.Println("Неверный формат числа сотрудников:", err)

			return
		}

		department := department.NewDepartment(countEmployees)

		for range countEmployees {
			if !scanner.Scan() {
				fmt.Println("Не удалось прочитать строку с требованием")

				return
			}

			line := strings.TrimSpace(scanner.Text())
			parts := strings.Fields(line)

			if len(parts) != partsCount {
				fmt.Println("Неверный формат строки:", line)

				return
			}

			operand := parts[0]

			temp, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Неверная температура:", parts[1], err)

				return
			}

			result := department.ProcessWorkerRequirement(operand, temp)
			fmt.Println(result)
		}
	}
}
