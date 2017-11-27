package main

import (
	"crypto/md5"
	"net/url"
	"runtime"
	"strconv"
	"time"

	//"compress/gzip"

	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/klauspost/pgzip"
	"github.com/remeh/sizedwaitgroup"
)

type Acesslog struct {
	Hash      string
	Type      string
	Time      string
	TTS       int
	SEIp      string
	Clientip  string
	Request   string
	Bytes     int
	Method    string
	Url       string
	Urlschema string
	Urlhost   string
	Urlpath   string
	Urlquery  string
	Mime      string
	Ua        string
}

type Ingestlogtest struct {
	Type         string
	Hash         string
	Time         string
	URL          string
	SEIp         string
	Urlschema    string
	Urlhost      string
	Urlpath      string
	Urlquery     string
	Urlfragment  string
	ServerIP     string
	BytesRead    int
	BytesToRead  int
	AssetSize    int
	Status       string
	IngestStatus string
}

type Ingestlog struct {
	Type             string
	Hash             string
	Time             string
	URL              string
	SEIp             string
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

// var wg sync.WaitGroup

var wg = sizedwaitgroup.New(200) //massimo numero di go routine per volta

func leggizip(file string) {
	defer wg.Done()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1) //esegue una go routine su tutti i processori -1

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
	for scan.Scan() {
		line := scan.Text()
		s := strings.Split(line, " ")
		if len(s) < 20 {
			continue
		}

		hasher := md5.New()
		hasher.Write([]byte(line))
		Hash := hex.EncodeToString(hasher.Sum(nil))
		val, err := client.SAdd("recordhashes", Hash).Result()
		// fmt.Println(val)
		// time.Sleep(3 * time.Second)
		if val == 0 { //se l'aggiunta dell'hash in redis è positiva prosegue altrimenti riprende il loop
			continue
		}
		t, err := time.Parse("02/Jan/2006:15:04:05", s[0][1:len(s[0])-7])
		if err != nil {
			fmt.Println(err)
		}
		Time := t.Format(time.RFC3339)
		//SEIp
		fileelements := strings.Split(file, "_")
		SEIp := fileelements[3]
		Type := fileelements[2]
		//gestiamo le url
		u, err := url.Parse(s[1])
		if err != nil {
			log.Fatal(err)
		}
		URL := s[1]
		Urlschema := u.Scheme
		Urlhost := u.Host
		Urlpath := u.Path
		Urlquery := u.RawQuery
		Urlfragment := u.Fragment
		//gestione url finita
		ServerIP := s[3]
		BytesRead, _ := strconv.Atoi(s[4])
		BytesToRead, _ := strconv.Atoi(s[5])
		AssetSize, _ := strconv.Atoi(s[6])
		Status := s[10]
		IngestStatus := s[15]
		record := &Ingestlogtest{Type: Type,
			Hash:         Hash,
			Time:         Time,
			URL:          URL,
			SEIp:         SEIp,
			Urlschema:    Urlschema,
			Urlhost:      Urlhost,
			Urlpath:      Urlpath,
			Urlquery:     Urlquery,
			Urlfragment:  Urlfragment,
			ServerIP:     ServerIP,
			BytesRead:    BytesRead,
			BytesToRead:  BytesToRead,
			AssetSize:    AssetSize,
			Status:       Status,
			IngestStatus: IngestStatus}

		out, err := json.Marshal(record)
		if err != nil {
			panic(err)
		}
		clienterr := client.LPush("codarecords", out).Err()
		if clienterr != nil {
			log.Fatal(clienterr)
		}
		//fmt.Printf("%+v\n", l)

	}
	return

}

func main() {
	for _, file := range os.Args[1:] {
		fmt.Println(file)
		wg.Add()
		go leggizip(file)
	}
	wg.Wait()

	return
}
