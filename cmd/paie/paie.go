package main

import (
	"assmat"
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const dbPath = "db.toml"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "paie: calcul des salaires d'une assistante maternelle\n")
	}
	flag.Parse()

	db := NewDb()
	toml.DecodeFile(dbPath, &db)

	cmds := map[string]cmd{
		"contract": &addCmd{db: &db},
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "error: need at least a command to execute\n")
		os.Exit(2)
	}

	err := cmds[args[0]].Run(args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	} else {
		f, err := os.Create(dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(2)
		}
		defer f.Close()
		enc := toml.NewEncoder(f)
		enc.Encode(db)
	}

	// contract: create/edit a contract
	// edit: edit a month sheet
	// show: show a month sheet
}

type cmd interface {
	Run(args []string) error
}

type db struct {
	Contracts map[string]assmat.Contract
}

func NewDb() db {
	return db{
		Contracts: make(map[string]assmat.Contract),
	}
}
