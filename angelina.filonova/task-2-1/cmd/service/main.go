package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultMinTemp = 15
	defaultMaxTemp = 30
	partsCount     = 2
)

type Department struct {
	minTemp  int
	maxTemp  int
	employee int
}

func NewDepartment(employee int) *Department {
	return &Department{
		minTemp:  defaultMinTemp,
		maxTemp:  defaultMaxTemp,
		employee: employee,
	}
}

func (d *Department) ProcessWorkerRequirement(operand string, temp int) int {
	switch operand {
	case ">=":
		if temp > d.minTemp {
			d.minTemp = temp
		}
	case "<=":
		if temp < d.maxTemp {
			d.maxTemp = temp
		}
	default:
		d.minTemp = temp
		d.maxTemp = temp
	}

	if d.minTemp > d.maxTemp {
		return -1
	}

	return d.minTemp
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		fmt.Println("Не удалось прочитать количество департаментов")

		return
	}

	countDepartmentsStr := strings.TrimSpace(scanner.Text())

	countDepartments, err := strconv.Atoi(countDepartmentsStr)
	if err != nil {
		fmt.Println("Неверный формат числа департаментов:", err)

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

		department := NewDepartment(countEmployees)

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
