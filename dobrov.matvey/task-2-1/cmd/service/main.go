package main

import (
	"fmt"
	"strconv"
)

func main() {
	var (
		countDepartments, countWorkers int
	)
	_, err := fmt.Scanln(&countDepartments)
	if err != nil {
		return
	}

	for i := 0; i < countDepartments; i++ {
		_, err = fmt.Scanln(&countWorkers)

		if err != nil {
			fmt.Println(-1)
			continue
		}

		minTemperature, maxTemperature := 15, 30
		broken := false
		for k := 0; k < countWorkers; k++ {
			if broken {
				var dump1, dump2 string
				fmt.Scanln(&dump1, &dump2)
				continue
			}

			var (
				needToIncrease                    bool
				desiredTemperature                int
				strOperand, strDesiredTemperature string
			)

			_, err = fmt.Scanln(&strOperand, &strDesiredTemperature)

			if err != nil {
				fmt.Println(-1)
				broken = true
				continue
			}

			if !parseDesiredTemperature(strOperand, strDesiredTemperature, &needToIncrease, &desiredTemperature) {
				fmt.Println(-1)
				broken = true
				continue
			}

			if needToIncrease {
				if desiredTemperature >= minTemperature {
					minTemperature = desiredTemperature
				}
			} else {
				if desiredTemperature <= maxTemperature {
					maxTemperature = desiredTemperature
				}
			}
			if minTemperature > maxTemperature {
				fmt.Println(-1)
				broken = true
				continue
			}
			fmt.Println(minTemperature)
		}
	}
}

func parseDesiredTemperature(strOperand string, strDesiredTemperature string, needToIncrease *bool, desiredTemperature *int) bool {

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
