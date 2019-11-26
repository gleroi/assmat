package main

import (
	"assmat"
	"fmt"
	"io"
	"os"
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

func writeSheet(w io.Writer, contract assmat.Contract, year int, month time.Month) {
	current := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	fmt.Fprintf(w, "#wd date       hour fee\n")
	for current.Month() <= month {
		if current.Weekday() == time.Monday {
			fmt.Fprintf(w, "\n")
		}
		dailyFee := 0.0
		if contract.WeekSchedule[current.Weekday()] > 0 {
			dailyFee = contract.DailyFee
		}
		fmt.Fprintf(w, "%s %s %.2f %.2f\n", weekDayStr[current.Weekday()],
			current.Format("02/01/2006"),
			contract.WeekSchedule[current.Weekday()],
			dailyFee)
		current = current.AddDate(0, 0, 1)
	}
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

	return nil
}
