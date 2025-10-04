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

	scanner.Scan()
	countDepartmentsStr := scanner.Text()
	countDepartments, _ := strconv.Atoi(strings.TrimSpace(countDepartmentsStr))

	for range countDepartments {
		scanner.Scan()
		countEmployeesStr := scanner.Text()
		countEmployees, _ := strconv.Atoi(strings.TrimSpace(countEmployeesStr))

		department := NewDepartment(countEmployees)

		for range countEmployees {
			scanner.Scan()
			line := strings.TrimSpace(scanner.Text())

			parts := strings.Fields(line)
			operand := parts[0]
			temp, _ := strconv.Atoi(parts[1])

			result := department.ProcessWorkerRequirement(operand, temp)
			fmt.Println(result)
		}
	}
}
