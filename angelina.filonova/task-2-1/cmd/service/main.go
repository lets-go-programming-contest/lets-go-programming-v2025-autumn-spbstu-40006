package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Department struct {
	minTemp  int
	maxTemp  int
	employee int
}

func NewDepartment(employee int) *Department {
	return &Department{
		minTemp:  15,
		maxTemp:  30,
		employee: employee,
	}
}

func (d *Department) ProcessWorkerRequirement(operand string, temp int) int {
	if operand == ">=" {
		if temp > d.minTemp {
			d.minTemp = temp
		}
	} else if operand == "<=" {
		if temp < d.maxTemp {
			d.maxTemp = temp
		}
	} else {
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

	for dep := 0; dep < countDepartments; dep++ {
		scanner.Scan()
		countEmployeesStr := scanner.Text()
		countEmployees, _ := strconv.Atoi(strings.TrimSpace(countEmployeesStr))

		department := NewDepartment(countEmployees)

		for emp := 0; emp < countEmployees; emp++ {
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
