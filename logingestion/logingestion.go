package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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
	RedirectedURL    string
	OSFailoverAction string
	BillingCookie    string
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

func fileingestion(file string) {

	scan := bufio.NewScanner(gr) //mettiamo tutto in un buffer che è rapido
	for scan.Scan() {
		line := scan.Text()
		s := strings.Split(line, " ") //splitta le linee secondo il delimitatore usato nel file di log, cambiare all'occorrenza

		if len(s) < 20 { // se i parametri sono meno di 20 allora ricomincia il loop, serve a evitare le linee che non ci interessano
			continue
		}

		hasher := md5.New()                                    //prepara a fare un hash
		hasher.Write([]byte(line))                             //hasha tutta la linea
		Hash := hex.EncodeToString(hasher.Sum(nil))            //estrae l'hash md5sum in versione quasi human readable
		val, err := client.SAdd("recordhashes", Hash).Result() //finalmente usiamo l'hash dentro a redis
		// fmt.Println(val)
		// time.Sleep(3 * time.Second)
		if val == 0 { //se l'aggiunta dell'hash in redis è positiva prosegue altrimenti riprende il loop
			continue //questo serve a ingestare solo cose nuove
		}
		t, err := time.Parse("[02/Jan/2006:15:04:05.000-0700]", s[0]) //quant'è bello parsare i timestamp in go :)
		if err != nil {
			fmt.Println(err)
		}
		Time := t.Format("2006-01-02T15:04:05.000Z") //ISO8601 mon amour

		//gestiamo le url secondo l'RFC ... non mi ricordo qual è
		u, err := url.Parse(s[1]) //prendi una URL, trattala male, falla a pezzi per ore...
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
		BytesRead, _ := strconv.Atoi(s[4])   //trasforma il valore in int
		BytesToRead, _ := strconv.Atoi(s[5]) //trasforma il valore in int
		AssetSize, _ := strconv.Atoi(s[6])   //trasforma il valore in int
		Status := s[10]
		IngestStatus := s[15]
		//creiamo un record con tutti i campi che ci interessano dentro
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

		out, err := json.Marshal(record) //con il record da solo ci facciamo una sega, serve encodarlo con il marshmellon in json
		if err != nil {
			panic(err)
		}
		clienterr := client.LPush("codarecords", out).Err() //metti il json dentro una lista di redis
		if clienterr != nil {
			log.Fatal(clienterr)
		}
	}
	//fmt.Printf("%+v\n", l)

}

func fileingestion2() {
	scan := bufio.NewScanner(gr) //mettiamo tutto in un buffer che è rapido
	for scan.Scan() {
		line := scan.Text()
		s := strings.Split(line, " ") //splitta le linee secondo il delimitatore usato nel file di log, cambiare all'occorrenza

		if len(s) < 20 { // se i parametri sono meno di 20 allora ricomincia il loop, serve a evitare le linee che non ci interessano
			continue
		}

		hasher := md5.New()                                    //prepara a fare un hash
		hasher.Write([]byte(line))                             //hasha tutta la linea
		Hash := hex.EncodeToString(hasher.Sum(nil))            //estrae l'hash md5sum in versione quasi human readable
		val, err := client.SAdd("recordhashes", Hash).Result() //finalmente usiamo l'hash dentro a redis
		// fmt.Println(val)
		// time.Sleep(3 * time.Second)
		if val == 0 { //se l'aggiunta dell'hash in redis è positiva prosegue altrimenti riprende il loop
			continue //questo serve a ingestare solo cose nuove
		}
		t, err := time.Parse("[02/Jan/2006:15:04:05.000-0700]", s[0]) //quant'è bello parsare i timestamp in go :)
		if err != nil {
			fmt.Println(err)
		}
		Time := t.Format("2006-01-02T15:04:05.000Z") //ISO8601 mon amour

		//gestiamo le url secondo l'RFC ... non mi ricordo qual è
		u, err := url.Parse(s[1]) //prendi una URL, trattala male, falla a pezzi per ore...
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
		BytesRead, _ := strconv.Atoi(s[4])   //trasforma il valore in int
		BytesToRead, _ := strconv.Atoi(s[5]) //trasforma il valore in int
		AssetSize, _ := strconv.Atoi(s[6])   //trasforma il valore in int
		Status := s[10]
		IngestStatus := s[15]
		//creiamo un record con tutti i campi che ci interessano dentro
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

		out, err := json.Marshal(record) //con il record da solo ci facciamo una sega, serve encodarlo con il marshmellon in json
		if err != nil {
			panic(err)
		}
		clienterr := client.LPush("codarecords", out).Err() //metti il json dentro una lista di redis
		if clienterr != nil {
			log.Fatal(clienterr)
		}
	}
	//fmt.Printf("%+v\n", l)

}
