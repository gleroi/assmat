package main

import (
	"assmat"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type addCmd struct {
	db *db
}

func (add addCmd) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("need a short name for the contract")
	}
	create(add.db, args[0])
	return nil
}

const contractPath = "/tmp/CONTRACT"

func create(globalDb *db, shortName string) error {
	// prepare a file contract
	db := NewDb()
	if contract, ok := globalDb.Contracts[shortName]; ok {
		db.Contracts[shortName] = contract
	} else {
		db.Contracts[shortName] = assmat.Contract{}
	}
	f := createContractFile()
	writeContract(f, db)
	f.Close()

	err := openEditor(contractPath)
	if err != nil {
		return err
	}

	// parse
	db = NewDb()
	_, err = toml.DecodeFile(contractPath, &db)
	if err != nil {
		return err
	}

	// validate
	for _, v := range db.Contracts {
		err := v.Validate()
		if err != nil {
			return err
		}
	}

	// print
	enc := toml.NewEncoder(os.Stdout)
	enc.Encode(db)

	// save
	mergeDb(globalDb, db)
	return nil
}

func mergeDb(dst *db, src db) {
	for k, v := range src.Contracts {
		dst.Contracts[k] = v
	}
}

func createContractFile() *os.File {
	f, err := os.Create(contractPath)
	if err != nil {
		panic(err)
	}
	return f
}

func writeContract(w io.Writer, db db) {
	enc := toml.NewEncoder(w)
	err := enc.Encode(db)
	if err != nil {
		panic(err)
	}
}
