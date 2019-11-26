package main

import (
	"assmat"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type monthCmd struct {
	db *db
}

var weekDayStr = []string{
	"Dim",
	"Lun",
	"Mar",
	"Mer",
	"Jeu",
	"Ven",
	"Sam",
}

type dayEntry struct {
	day        time.Time
	WorkedHour float64
	dailyFee   float64
}

func newDayEntry(day time.Time, contract assmat.Contract) dayEntry {
	de := dayEntry{
		day:        time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC),
		WorkedHour: contract.WeekSchedule[day.Weekday()],
	}
	if de.WorkedHour > 0 {
		de.dailyFee = contract.DailyFee
	}
	return de
}

func (de dayEntry) String() string {
	return fmt.Sprintf("%s %s %.2f %.2f", weekDayStr[de.day.Weekday()],
		de.day.Format("02/01/2006"), de.WorkedHour, de.dailyFee)
}

func (de *dayEntry) UnmarshalText(text []byte) error {
	var dayStr string
	var nameStr string

	_, err := fmt.Sscanf(string(text), "%s %s %f %f", &nameStr,
		&dayStr, &de.WorkedHour, &de.dailyFee)
	if err != nil {
		return err
	}
	de.day, err = time.Parse("02/01/2006", dayStr)
	return err
}

func writeSheet(w io.Writer, contract assmat.Contract, year int, month time.Month) {
	current := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	fmt.Fprintf(w, "#wd date       hour fee\n")
	for current.Month() <= month {
		if current.Weekday() == time.Monday {
			fmt.Fprintf(w, "\n")
		}
		de := newDayEntry(current, contract)
		fmt.Fprintf(w, "%s\n", de)
		current = current.AddDate(0, 0, 1)
	}
}

func readSheet(r io.Reader, contract assmat.Contract, year int, month time.Month) error {
	scan := bufio.NewScanner(r)
	count := 0

	feeSum := 0.0
	for scan.Scan() {
		count++
		if scan.Text() == "" || strings.HasPrefix(scan.Text(), "#") {
			continue
		}
		var de dayEntry
		err := de.UnmarshalText(scan.Bytes())
		if err != nil {
			return fmt.Errorf("line %d : %w", count, err)
		}
		feeSum += de.dailyFee
	}
	fmt.Fprintf(os.Stdout, "Basic salary: %.2f €\n", contract.BaseSalary())
	fmt.Fprintf(os.Stdout, "Daily fees: %.2f €\n", feeSum)
	fmt.Fprintf(os.Stdout, "Total: %.2f €\n", float64(contract.BaseSalary())+feeSum)
	return nil
}

func (cmd *monthCmd) Run(args []string) error {
	now := time.Now()

	if len(args) < 1 {
		return fmt.Errorf("need a short name for the contract")
	}
	contract, ok := cmd.db.Contracts[args[0]]
	if !ok {
		return fmt.Errorf("unknown contract name: %s", args[0])
	}

	f, err := os.Create("/tmp/MONTH_SHEET")
	if err != nil {
		return err
	}
	writeSheet(f, contract, now.Year(), now.Month())
	f.Close()

	err = openEditor("/tmp/MONTH_SHEET")
	if err != nil {
		return err
	}

	//TODO: save month sheet

	//TODO: compute total month salary
	f, err = os.Open("/tmp/MONTH_SHEET")
	if err != nil {
		return err
	}
	return readSheet(f, contract, now.Year(), now.Month())
}
