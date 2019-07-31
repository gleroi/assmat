package assmat

import "math"

// DaysInWeek is the number of day in a week
const DaysInWeek = 7
const WeeksInYear = 52
const MonthsInYear = 12
const PajEmploiNetDailySalaryLimit = 39.12

// WeekSchedule is the number of hours worked each day in a week
type WeekSchedule [DaysInWeek]float64

// Hours return the number of hours worked in a full week
func (s WeekSchedule) Hours() float64 {
	var weekHours float64
	for _, dayHour := range s {
		weekHours += dayHour
	}
	return weekHours
}

func (s WeekSchedule) Days() float64 {
	var days float64
	for _, d := range s {
		if d > 0 {
			days += 1.0
		}
	}
	return days
}

type Salary float64

func (s Salary) Round() float64 {
	i := math.Round(float64(s) * 100)
	return float64(i) / 100
}

// Contract regroup all to information to calculate the monthly salary of assmat
type Contract struct {
	Name              string
	HourlyRate        float64
	DailyFee          float64
	WeekSchedule      WeekSchedule
	WorkedWeeksInYear int
}

// MonthlySalary computes the base monthly salary
func (c Contract) BaseSalary() Salary {
	return Salary(c.WeekSchedule.Hours() * c.HourlyRate * float64(c.WorkedWeeksInYear) / MonthsInYear)
}

// WorkedHours is the monthly hours count worked in a month to declare
func (c Contract) WorkedHours() float64 {
	return c.WeekSchedule.Hours() * float64(c.WorkedWeeksInYear) / MonthsInYear
}

// WorkedHours is the monthly days count worked in a month to declare
func (c Contract) WorkedDays() float64 {
	return c.WeekSchedule.Days() * float64(c.WorkedWeeksInYear) / MonthsInYear
}

func (c Contract) DailySalary() Salary {
	return Salary(float64(c.BaseSalary()) / c.WorkedDays())
}

func (c Contract) Validate() []error {
	var errors []error

	return errors
}
