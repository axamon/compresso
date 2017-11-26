package main

import (
	"net/url"
	"time"

	//"compress/gzip"

	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/klauspost/pgzip"
)

type Acesslog struct {
	Time     string
	TTS      int
	clientip string
	request  string
	bytes    int
	method   string
	url      string
	mime     string
	ua       string
	unused1  string
	unused2  string
	unused3  string
	unused4  string
	unused5  string
}

type Ingestlog struct {
	Time             string
	URL              string
	Urlschema        string
	Urlhost          string
	Urlpath          string
	Urlquery         string
	Urlfragment      string
	FailOverSvrList  string
	ServerIP         string
	BytesRead        int
	BytesToRead      int
	AssetSize        int
	DownloadComplete string
	DownloadTime     string
	ReadCallBack     string
	Status           string
	Mime             string
	Revaldidation    string
	CDSDomain        string
	ConnectionInfo   string
	IngestStatus     string
	RedirectedUrl    string
	OSFailoverAction string
	BillingCookie    string
}

var wg sync.WaitGroup

func leggizip(file string, wg *sync.WaitGroup) {
	defer wg.Done()
	//runtime.GOMAXPROCS(1)
	// runtime.NumCPU()
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// gr, err := gzip.NewReader(f)
	gr, err := pgzip.NewReaderN(f, 4096, 100)

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
	return
}

func leggizip2(file string, wg *sync.WaitGroup) {
	defer wg.Done()
	//runtime.GOMAXPROCS(1)
	// runtime.NumCPU()
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// gr, err := gzip.NewReader(f)
	gr, err := pgzip.NewReaderN(f, 4096, 100)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	scan := bufio.NewScanner(gr)
	l := Ingestlog{}
	for scan.Scan() {
		line := scan.Text()
		s := strings.Split(line, " ")
		if len(s) < 20 {
			continue
		}
		//gestiamo il tempo
		t, err := time.Parse("02/Jan/2006:15:04:05", s[0][1:len(s[0])-7])
		if err != nil {
			fmt.Println(err)
		}
		l.Time = t.Format(time.RFC3339)
		//gestione url finita
		//gestiamo le url
		u, err := url.Parse(s[1])
		if err != nil {
			log.Fatal(err)
		}
		l.URL = s[1]
		l.Urlschema = u.Scheme
		l.Urlhost = u.Host
		l.Urlpath = u.Path
		l.Urlquery = u.RawQuery
		l.Urlfragment = u.Fragment
		//gestione url finita
		l.ServerIP = s[3]
		fmt.Printf("%#v\n", l)

	}
	return

}

func main() {
	for _, file := range os.Args[1:] {
		fmt.Println(file)
		wg.Add(1)
		// go leggizipEvolved(file)
		go leggizip2(file, &wg)
		//Leggi(file)
	}
	wg.Wait()

	return
}
