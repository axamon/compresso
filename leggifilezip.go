package main

import (
	"bufio"
	"strings"

	gzip "github.com/klauspost/pgzip"
	//"compress/gzip"
	"fmt"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

func leggizip(file string) {
	// runtime.GOMAXPROCS(2)
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
	scanner := bufio.NewScanner(gr)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		_ = strings.Split(line, " ")
		//fmt.Println(s)
		// strchan <- scanner.Text()
	}

	wg.Done()
	return
}

func main() {
	for _, file := range os.Args[:] {
		fmt.Println(file)
		wg.Add(1)
		// go leggizipEvolved(file)
		leggizip(file)
	}
	wg.Wait()
	return
}
