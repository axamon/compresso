package leggizipevolved

import (
	"io"

	"github.com/klauspost/pgzip"

	//"compress/gzip"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func LeggizipEvolved(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gr, err := pgzip.NewReaderN(f, 4096, 10000)
	if err != nil {
		log.Fatal(err)
	}
	defer gr.Close()

	cr := csv.NewReader(gr)

	cr.Comma = ' '          //specifica il delimitatore dei campi
	cr.FieldsPerRecord = -1 //accetta numero di campi variabili
	cr.Comment = '#'
	//cr.Comma = delimiter //specifica il delimitatore dei campi
	cr.LazyQuotes = true
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(rec)
	}
	return
}
