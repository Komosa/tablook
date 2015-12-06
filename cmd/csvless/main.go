package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/komosa/tablook"
)

func process(in io.Reader) error {
	r := csv.NewReader(in)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	tbl, err := tablook.New(records)
	if err != nil {
		return err
	}

	err = tbl.Show()
	return err
}

func main() {
	err := processFileOrStdin()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func processFileOrStdin() error {
	if len(os.Args) != 2 {
		return process(os.Stdin)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer f.Close()
	return process(f)
}
