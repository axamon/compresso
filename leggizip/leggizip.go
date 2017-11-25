package leggizip

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
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
		continue
	}
	defer gr.Close()

	// scanner := bufio.NewScanner(gr)

	// 		scanner.Split(bufio.ScanLines)

	// 		for scanner.Scan() {
	// 			fmt.Println(scanner.Text())
	// 			strchan <- scanner.Text()
	// 		}

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
