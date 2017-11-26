package main

import (

	//"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/klauspost/pgzip"
)

var wg sync.WaitGroup

func leggizip(file string) error {
	runtime.GOMAXPROCS(1)
	// runtime.NumCPU()
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//gr, err := gzip.NewReader(f)
	gr, err := pgzip.NewReaderN(f, 4096, 10000)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
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
	wg.Done()
	return nil
}

func main() {
	for _, file := range os.Args[1:] {
		fmt.Println(file)
		wg.Add(1)
		// go leggizipEvolved(file)
		go leggizip(file)
		//Leggi(file)
	}
	wg.Wait()
	return
}
