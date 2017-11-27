package main

import (
	"crypto/md5"
	"net/url"
	"runtime"
	"strconv"
	"time"

	//"compress/gzip"

	"bufio"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/go-redis/redis"
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
	Hash             string
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

func leggizip(file string, wg *sync.WaitGroup) {
	defer wg.Done()
	runtime.GOMAXPROCS(runtime.NumCPU())      //esegue una go routine su tutti i processori
	client := redis.NewClient(&redis.Options{ //connettiti a Redis server
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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

		hasher := md5.New()
		hasher.Write([]byte(line))
		l.Hash = hex.EncodeToString(hasher.Sum(nil))
		val, err := client.SAdd("recordhashes", l.Hash).Result()
		// fmt.Println(val)
		// time.Sleep(3 * time.Second)
		if val == 0 { //se l'aggiunta dell'hash in redis è positiva prosegue altrimenti riprende il loop
			continue
		}
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
		l.BytesRead, _ = strconv.Atoi(s[4])
		l.BytesToRead, _ = strconv.Atoi(s[5])
		l.AssetSize, _ = strconv.Atoi(s[6])
		l.Status = s[10]
		l.IngestStatus = s[15]
		err := client.LPush("codarecords", l).Err()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%#v\n", l)

	}
	return

}

func main() {
	for _, file := range os.Args[1:] {
		fmt.Println(file)
		wg.Add(1)
		// go leggizipEvolved(file)
		go leggizip(file, &wg)
		//Leggi(file)
	}
	wg.Wait()

	return
}
