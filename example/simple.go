package main

import (
	"fmt"
	"strings"

	"github.com/komosa/tablook"
)

func main() {
	d := [][]string{
		{"hdr1", "hdr2", "hdr3", "header4"},
		{"1", "2", "πœę©þąśðł", "test cztery"},
		{"multiple-cell characters: ", "つのだ☆HIRO", "c,", "d"},
		{"3 this", "row", "will be", "very lo" + strings.Repeat("o", 100) + "ng"},
	}
	for i := 4; i < 100; i++ {
		s := fmt.Sprint(i)
		ss := []string{s, s, s, s}
		d = append(d, ss)
	}

	t, err := tablook.New(d)
	if err != nil {
		panic(err)
	}
	t.Show()
}
