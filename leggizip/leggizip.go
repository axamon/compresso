package leggizip

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func Leggizip(file string) {
	f, err := os.Open(file)
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
