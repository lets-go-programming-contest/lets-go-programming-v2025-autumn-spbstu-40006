package department

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
