package main

import (
	"errors"
	"fmt"
	"strconv"
)

func main() {
	var deptCount, workerCount int

	if _, err := fmt.Scanln(&deptCount); err != nil {
		return
	}

	for i := 0; i < deptCount; i++ {
		if _, err := fmt.Scanln(&workerCount); err != nil {
			fmt.Println(-1)
			continue
		}
		handleDepartment(workerCount)
	}
}

func handleDepartment(workers int) {
	minT, maxT := 15, 30
	failed := false

	for i := 0; i < workers; i++ {
		if failed {
			var skipA, skipB string
			if _, err := fmt.Scanln(&skipA, &skipB); err != nil {
				fmt.Println(-1)
				continue
			}
			fmt.Println(-1)
			continue
		}

		var (
			opSymbol string
			tempStr  string
		)

		if _, err := fmt.Scanln(&opSymbol, &tempStr); err != nil {
			fmt.Println(-1)
			failed = true
			continue
		}

		var (
			shouldRaise bool
			targetTemp  int
		)

		if err := interpretConstraint(opSymbol, tempStr, &shouldRaise, &targetTemp); err != nil {
			fmt.Println(-1)
			failed = true
			continue
		}

		if !updateTemperatureBounds(shouldRaise, targetTemp, &minT, &maxT) {
			fmt.Println(-1)
			failed = true
			continue
		}

		fmt.Println(minT)
	}
}

func updateTemperatureBounds(increase bool, desired int, minT, maxT *int) bool {
	if increase {
		if desired > *minT {
			*minT = desired
		}
	} else {
		if desired < *maxT {
			*maxT = desired
		}
	}
	return *minT <= *maxT
}

var errBadOperator = errors.New("invalid operator")

func interpretConstraint(
	op string,
	val string,
	shouldIncrease *bool,
	target *int,
) error {
	switch op {
	case ">=":
		*shouldIncrease = true
	case "<=":
		*shouldIncrease = false
	default:
		return fmt.Errorf("%w: %s", errBadOperator, op)
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("invalid value %q: %w", val, err)
	}

	*target = num
	return nil
}
