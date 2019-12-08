package assmat

import (
	"fmt"
	"math"
	"testing"
)

var contracts = []Contract{
	Contract{
		HourlyRate:        3.62,
		DailyFee:          3.50,
		WorkedWeeksInYear: 45,
		WeekSchedule:      WeekSchedule{9, 9, 9, 9, 9, 0, 0},
	},
	Contract{
		HourlyRate:        5,
		DailyFee:          3.10,
		WorkedWeeksInYear: 44,
		WeekSchedule:      WeekSchedule{0, 0, 4.5, 0, 0, 0, 0},
	},
	Contract{
		HourlyRate:        3.90,
		DailyFee:          4.00,
		WorkedWeeksInYear: 45,
		WeekSchedule:      WeekSchedule{9, 9, 0, 9, 9, 0, 0},
	},
}

func TestContractBaseSalary(t *testing.T) {
	expected := []float64{
		610.88,
		82.50,
		526.50,
	}

	for i, contract := range contracts {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			salary := contract.BaseSalary()
			if salary.Round() != expected[i] {
				t.Errorf("expected %f, got %f", expected[i], salary)
			}
		})
	}
}

func TestContractWorkedHours(t *testing.T) {
	expected := []float64{
		168.75,
		16.5,
		135,
	}

	for i, contract := range contracts {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			hours := contract.WorkedHours()
			if hours != expected[i] {
				t.Errorf("expected %f, got %f", expected[i], hours)
			}
		})
	}
}

func TestContractWorkedDays(t *testing.T) {
	expected := []float64{
		19,
		4,
		15,
	}

	for i, contract := range contracts {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			days := math.Round(contract.WorkedDays())
			if days != expected[i] {
				t.Errorf("expected %f, got %f", expected[i], days)
			}
		})
	}
}
