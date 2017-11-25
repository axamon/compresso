package main

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("data.csv.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	defer gr.Close()

	cr := csv.NewReader(gr)

	cr.Comma = ';' //specifica il delimitatore dei campi
	rec, err := cr.Read()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range rec {
		fmt.Printf(v + " ")
	}
}
