package main

import (
	"assmat"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
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
	mealFee    float64
}

func newDayEntry(day time.Time, contract assmat.Contract) dayEntry {
	de := dayEntry{
		day:        time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC),
		WorkedHour: contract.WeekSchedule[day.Weekday()],
	}
	if de.WorkedHour > 0 {
		de.dailyFee = contract.DailyFee
		de.mealFee = contract.MealFee
	}
	return de
}

func (de dayEntry) String() string {
	return fmt.Sprintf("%s\t%s\t%.2f\t%.2f\t%.2f", weekDayStr[de.day.Weekday()],
		de.day.Format("02/01/2006"), de.WorkedHour, de.dailyFee, de.mealFee)
}

func (de *dayEntry) UnmarshalText(text []byte) error {
	var dayStr string
	var nameStr string

	_, err := fmt.Sscanf(string(text), "%s %s %f %f %f", &nameStr,
		&dayStr, &de.WorkedHour, &de.dailyFee, &de.mealFee)
	if err != nil {
		return err
	}
	de.day, err = time.Parse("02/01/2006", dayStr)
	return err
}

func writeNewSheet(w io.Writer, contract assmat.Contract, year int, month time.Month) {
	current := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	entries := make([]dayEntry, 0, 31)
	for current.Month() <= month {
		de := newDayEntry(current, contract)
		entries = append(entries, de)
		current = current.AddDate(0, 0, 1)
	}
	writeSheet(w, contract, entries)
}

func writeSheet(iw io.Writer, contract assmat.Contract, entries []dayEntry) {
	w := tabwriter.NewWriter(iw, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "#wd\tdate\thour\tfee\tmeal\n")
	for _, de := range entries {
		current := de.day
		if current.Weekday() == time.Monday {
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "%s\n", de)
	}
	fmt.Fprintf(w, "\n")
	writeSheetSummary(w, contract, entries)
	w.Flush()
}

func writeSheetSummary(w io.Writer, contract assmat.Contract, entries []dayEntry) {
	feeSum := 0.0
	for _, de := range entries {
		feeSum += de.dailyFee
	}

	fmt.Fprintf(w, "# Basic salary:\t %.2f €\n", contract.BaseSalary())
	fmt.Fprintf(w, "# Daily fees:\t %.2f €\n", feeSum)
	fmt.Fprintf(w, "# Total:\t %.2f €\n", float64(contract.BaseSalary())+feeSum)
}

func readSheet(r io.Reader, contract assmat.Contract) ([]dayEntry, error) {
	scan := bufio.NewScanner(r)
	count := 0
	entries := make([]dayEntry, 0, 31)
	for scan.Scan() {
		count++
		if scan.Text() == "" || strings.HasPrefix(scan.Text(), "#") {
			continue
		}
		var de dayEntry
		err := de.UnmarshalText(scan.Bytes())
		if err != nil {
			return entries, fmt.Errorf("line %d : %w", count, err)
		}
		entries = append(entries, de)
	}

	return entries, nil
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
	writeNewSheet(f, contract, now.Year(), now.Month())
	f.Close()

	err = openEditor("/tmp/MONTH_SHEET")
	if err != nil {
		return err
	}

	// read sheet
	f, err = os.Open("/tmp/MONTH_SHEET")
	if err != nil {
		return err
	}
	entries, err := readSheet(f, contract)
	if err != nil {
		return err
	}

	//TODO: save month sheet
	//TODO: compute total month salary
	f, err = os.Create(fmt.Sprintf("%04d_%02d_%s.txt", now.Year(), now.Month(), args[0]))
	if err != nil {
		return err
	}
	defer f.Close()
	writeSheet(f, contract, entries)
	writeSheetSummary(os.Stdout, contract, entries)
	return nil
}
