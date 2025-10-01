package main

import (
	"fmt"
	"strconv"
)

func main() {
	var (
		countDepartments int
		countWorkers     int
	)

	_, err := fmt.Scanln(&countDepartments)
	if err != nil {
		return
	}

	for range countDepartments {
		_, err = fmt.Scanln(&countWorkers)
		if err != nil {
			fmt.Println(-1)

			continue
		}

		processDepartment(countWorkers)
	}
}

func processDepartment(countWorkers int) {
	minTemperature, maxTemperature := 15, 30
	broken := false

	for range countWorkers {
		if broken {
			var dump1, dump2 string
			if _, err := fmt.Scanln(&dump1, &dump2); err != nil {
				continue
			}

			continue
		}

		var (
			needIncrease bool
			desired      int
			operand      string
			desiredStr   string
		)

		if _, err := fmt.Scanln(&operand, &desiredStr); err != nil {
			fmt.Println(-1)
			broken = true

			continue
		}

		// wsl: if с проверкой — отдельным блоком после присваиваний/ввода
		if !parseDesiredTemperature(operand, desiredStr, &needIncrease, &desired) {
			fmt.Println(-1)
			broken = true

			continue
		}

		if !applyConstraint(needIncrease, desired, &minTemperature, &maxTemperature) {
			fmt.Println(-1)
			broken = true

			continue
		}

		fmt.Println(minTemperature)
	}
}

func applyConstraint(needUp bool, desired int, minT *int, maxT *int) bool {
	if needUp {
		if desired >= *minT {
			*minT = desired
		}
	} else {
		if desired <= *maxT {
			*maxT = desired
		}
	}

	return *minT <= *maxT
}

func parseDesiredTemperature(
	strOperand string,
	strDesiredTemperature string,
	needToIncrease *bool,
	desiredTemperature *int,
) bool {
	switch strOperand {
	case ">=":
		*needToIncrease = true
	case "<=":
		*needToIncrease = false
	default:
		return false
	}

	value, err := strconv.Atoi(strDesiredTemperature)
	if err != nil {
		return false
	}

	*desiredTemperature = value

	return true
}
