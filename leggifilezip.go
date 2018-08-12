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

//Accesslog tipo per transaction
type Accesslog struct {
	Hash      string
	Type      string
	Time      string
	TTS       int
	SEIp      string
	Clientip  string
	Request   string
	Bytes     int
	Method    string
	URL       string
	Urlschema string
	Urlhost   string
	Urlpath   string
	Urlquery  string
	Mime      string
	Ua        string
}

//Ingestlogtest recors dei log ingestion
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

//Ingestlog fields per Ingestion logs
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

// var wg sync.WaitGroup

var wg = sizedwaitgroup.New(200) //massimo numero di go routine per volta
//var Test string = "pippo"

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
	gr, err := pgzip.NewReaderN(f, 4096, 100) //sfrutta il gzip con steroide che legge nel futuro per andare più veloce assai

	if err != nil { //se però si impippa qualcosa allora blocca tutto
		log.Fatal(err)
		os.Exit(1)
	}

	fileelements := strings.Split(file, "_") //prende il nome del file di log e recupera i campi utili
	Type := fileelements[1]                  //qui prede il tipo di log
	SEIp := fileelements[3]                  //qui prende l'ip della cache

	if Type == "accesslog" { //se il tipo di log è "accesslog" allora fa qualcosa che ancora non ho finito di fare
		scan := bufio.NewScanner(gr)
		for scan.Scan() {
			line := scan.Text()
			s := strings.Split(line, "\t")
			// fmt.Println("dopo")
			t, err := time.Parse("[02/Jan/2006:15:04:05.000-0700]", s[0]) //converte i timestamp come piacciono a me
			if err != nil {
				fmt.Println(err)
			}
			Time := t.Format("2006-01-02T15:04:05.000Z") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
			fmt.Println(Time)
			continue
		}
	}

	if Type == "ingestlog" {
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
	return //terminata la Go routine!!! :)
}

func leggizip2(file string) {
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
	gr, err := pgzip.NewReaderN(f, 4096, 100) //sfrutta il gzip con steroide che legge nel futuro per andare più veloce assai

	if err != nil { //se però si impippa qualcosa allora blocca tutto
		log.Fatal(err)
		os.Exit(1)
	}

	fileelements := strings.Split(file, "_") //prende il nome del file di log e recupera i campi utili
	Type := fileelements[1]                  //qui prede il tipo di log
	SEIp := fileelements[3]                  //qui prende l'ip della cache

	if Type == "accesslog" { //se il tipo di log è "accesslog" allora fa qualcosa che ancora non ho finito di fare
		scan := bufio.NewScanner(gr)
		for scan.Scan() {
			line := scan.Text()
			if !strings.HasPrefix(line, "[") { //se la linea non inzia con [ allora salta
				continue
			}
			s := strings.Split(line, "\t")
			//t, err := time.Parse("[02/Jan/2006:15:04:05.000-0700]", s[0]) //converte i timestamp come piacciono a me
			if err != nil {
				fmt.Println(err)
			}
			//Time := t.Unix()
			//Time := t.Format("2006-01-02T15:04:05.000Z") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
			//fmt.Println(Time)
			var speed, tts, bytes float64

			tts, err = strconv.ParseFloat(s[1], 8)
			if err != nil {
				log.Fatal(err.Error())
			}

			bytes, err := strconv.ParseFloat(s[4], 8)
			if err != nil {
				log.Fatal(err.Error())
			}

			speed = (bytes / tts)
			clientip := s[2]
			status := s[3]
			ua := s[8]

			u, err := url.Parse(s[6]) //prendi una URL, trattala male, falla a pezzi per ore...
			if err != nil {
				log.Fatal(err)
			}
			//Urlschema := u.Scheme
			//Urlhost := u.Host
			Urlpath := u.Path
			//fmt.Println(Urlpath)
			//Urlquery := u.RawQuery
			//Urlfragment := u.Fragment
			if strings.Contains(Urlpath, "videoteca") {
				pezziurl := strings.Split(Urlpath, "/")
				//fmt.Println(pezziurl)
				idvideoteca := pezziurl[6]
				encoding := pezziurl[10]
				hasher := md5.New()                                  //prepara a fare un hash
				hasher.Write([]byte(clientip + idvideoteca + ua))    //hasha tutta la linea
				Hash := hex.EncodeToString(hasher.Sum(nil))          //estrae l'hash md5sum in versione quasi human readable
				_, err := client.SAdd("recordhashes", Hash).Result() //finalmente usiamo l'hash dentro a redis
				if err != nil {
					log.Fatal(err.Error())
				}
				//fmt.Println(idvideoteca)
				if speed < 0.87 {
					fmt.Printf("%v %v %.3f %v %v %v\n", Hash, idvideoteca, speed, status, clientip, encoding)
				}
			}

		}
	}

	if Type == "ingestlog" {
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
	return //terminata la Go routine!!! :)
}

func main() {
	for _, file := range os.Args[1:] {
		fmt.Println(file)
		wg.Add()
		go leggizip2(file)
	}
	wg.Wait()

	return
}
