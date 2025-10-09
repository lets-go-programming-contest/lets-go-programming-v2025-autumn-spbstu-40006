package main

import "fmt"

func main() {
	var (
		department_amount, employee_amount, temperature    int
		cur_temp_limit, lower_temp_limit, upper_temp_limit int
		str                                                string
		err                                                error
	)

	_, err = fmt.Scan(&department_amount)
	if err != nil || department_amount < 1 || department_amount > 1000 {
		fmt.Println(-1)
		return
	}

	for i := 0; i < department_amount; i++ {
		temperature = 0
		lower_temp_limit = 15
		upper_temp_limit = 30

		_, err = fmt.Scan(&employee_amount)
		if err != nil || employee_amount < 1 || employee_amount > 1000 {
			fmt.Println(-1)
			return
		}

		for j := 0; j < employee_amount; j++ {
			_, err = fmt.Scan(&str, &cur_temp_limit)
			if err != nil {
				fmt.Println(-1)
				return
			}

			if str == ">=" && cur_temp_limit <= upper_temp_limit {
				lower_temp_limit = max(lower_temp_limit, cur_temp_limit)
				temperature = max(lower_temp_limit, temperature)

			} else if str == "<=" && cur_temp_limit >= lower_temp_limit {
				upper_temp_limit = min(upper_temp_limit, cur_temp_limit)
				temperature = min(upper_temp_limit, temperature)

			} else {
				fmt.Println(-1)
				return
			}

			fmt.Println(temperature)
		}
	}
}
