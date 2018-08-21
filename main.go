package main

import (
	"crypto/md5"
	"net/url"
	"runtime"
	"sort"
	"strconv"

	"gonum.org/v1/gonum/stat"

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

const (
	bitstoMB = 0.000000125
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

var wg = sizedwaitgroup.New(200) //massimo numero di go routine per volta

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
		log.Fatal(err.Error())
		os.Exit(1)
	}

	fileelements := strings.Split(file, "_") //prende il nome del file di log e recupera i campi utili
	Type := fileelements[1]                  //qui prede il tipo di log
	//SEIp := fileelements[3]                  //qui prende l'ip della cache

	if Type == "accesslog" { //se il tipo di log è "accesslog" allora fa qualcosa che ancora non ho finito di fare
		scan := bufio.NewScanner(gr)
		for scan.Scan() {
			line := scan.Text()
			err := escludidoppioni(line)
			if err != nil {
				continue
			}
			if !strings.HasPrefix(line, "[") { //se la linea non inzia con [ allora salta
				continue
			}
			s := strings.Split(line, "\t")
			u, err := url.Parse(s[6]) //parsa la URL nelle sue componenti
			if err != nil {
				log.Fatal(err)
			}
			/* Urlschema := u.Scheme
			if Urlschema != "https" { //fa passare solo le URL richieste via WEB
				continue
			} */

			//t, err := time.Parse("[02/Jan/2006:15:04:05.000-0700]", s[0]) //converte i timestamp come piacciono a me
			if err != nil {
				fmt.Println(err)
			}
			//	Time := t.Unix()
			//Time := t.Format("2006-01-02T15:04:05.000Z") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
			//fmt.Println(Time)
			var speed, tts, bytes float64

			tts, err = strconv.ParseFloat(s[1], 8)
			if err != nil {
				log.Fatal(err.Error())
			}

			bytes, err = strconv.ParseFloat(s[4], 8)
			if err != nil {
				log.Fatal(err.Error())
			}

			speed = (bytes / tts)
			clientip := s[2]
			//status := s[3]
			ua := s[8]

			//fmt.Println(Urlschema)
			//Urlhost := u.Host
			Urlpath := u.Path
			//fmt.Println(Urlpath)
			//Urlquery := u.RawQuery
			//Urlfragment := u.Fragment
			pezziurl := strings.Split(Urlpath, "/")
			//fmt.Println(pezziurl)
			/* if len(pezziurl) < 11 {
				continue
			} */
			if ok := strings.HasPrefix(Urlpath, "videoteca"); ok == true { //Prende solo i chunk video per
				//fmt.Println(pezziurl)
				//fmt.Println(Urlpath)
				idvideoteca := pezziurl[6]
				//fmt.Println(idvideoteca)
				//encoding := pezziurl[10]
				//fmt.Println(encoding)
				//re := regexp.MustCompile(`QualityLevels\(([0-9]+)\)$`)
				//bitratestr := re.FindStringSubmatch(encoding)[1]
				//bitrate, _ := strconv.ParseFloat(bitratestr, 8)
				if err != nil {
					log.Fatal(err.Error())
				}
				//bitrateMB := bitrate * bitstoMB

				hasher := md5.New()                                 //prepara a fare un hash
				hasher.Write([]byte(clientip + idvideoteca + ua))   //hasha tutta la linea
				Hash := hex.EncodeToString(hasher.Sum(nil))         //estrae l'hash md5sum in versione quasi human readable
				_, err = client.SAdd("recordhashes", Hash).Result() //finalmente usiamo l'hash dentro a redis
				if err != nil {
					log.Fatal(err.Error())
				}
				//fmt.Println(idvideoteca)

				//	fmt.Printf("%v %v %v %.3f %v %v %.3f %.3f\n", Time, Hash, idvideoteca, speed, status, clientip, bitrateMB, speed-bitrateMB)
				/* clientiparsed := net.ParseIP(clientip)
				//fmt.Println(clientiparsed)
				clientipint := IPv4ToInt(clientiparsed)
				idvideotecaint, err := strconv.Atoi(idvideoteca)
				if err != nil {
					log.Fatal(err.Error())
				} */
				ingestafruizioni(Hash, clientip, idvideoteca, speed)
				//hm, _ := stats.HarmonicMean([]float64{1, 2, 3, 4, 5})
			}
			if ok := strings.Contains(Urlpath, "DASH"); ok == true { //Prende solo i chunk DASH
				idvideoteca := pezziurl[6]
				hasher := md5.New()                                 //prepara a fare un hash
				hasher.Write([]byte(clientip + idvideoteca + ua))   //hasha tutta la linea
				Hash := hex.EncodeToString(hasher.Sum(nil))         //estrae l'hash md5sum in versione quasi human readable
				_, err = client.SAdd("recordhashes", Hash).Result() //finalmente usiamo l'hash dentro a redis
				if err != nil {
					log.Fatal(err.Error())
				}
				//	fmt.Printf("%v %v %v %.3f %v %v\n", Time, Hash, idvideoteca, speed, status, clientip)
				/* clientiparsed := net.ParseIP(clientip)
				clientipint := IPv4ToInt(clientiparsed)
				idvideotecaint, err := strconv.Atoi(idvideoteca)
				if err != nil {
					log.Fatal(err.Error())
				} */
				ingestafruizioni(Hash, clientip, idvideoteca, speed)
			}
		}
	}

	return //terminata la Go routine!!! :)
}

type Fruizioni struct {
	Hashfruizione map[string]bool
	Clientip      map[string]string
	Idvideoteca   map[string]string
	Details       map[string][]float64 `json:"-"`
}

type Fruizioniexport struct {
	Hashfruizione string `json:"-"` //Non permette a json di esportare il campo
	Clientip      string
	Idvideoteca   string
	Errori        int
}

//F contiene tutte le informazioni delle varie fruizioni
//è la variabile che verrà resa persistente su disco
var F = Fruizioni{}

func main() {

	hashline = make(map[string]bool)

	//se il file hashlinefile non esiste lo crea
	//gobfile è il file dove verrà resa persistente
	if _, err := os.Stat(hashlinefile); os.IsNotExist(err) {
		os.Create(hashlinefile)
	}

	//Carico dentro hashline i dati salvati precedentemente
	err := load(hashlinefile, &hashline) //Carica in Contatori i dati salvati sul gobfile
	if err != nil {
		fmt.Println(err.Error()) //se da errore forse manca il file... non importa se è il primo avvio
	}
	defer save(hashlinefile, hashline)

	F.Hashfruizione = make(map[string]bool)
	F.Clientip = make(map[string]string)
	F.Idvideoteca = make(map[string]string)
	F.Details = make(map[string][]float64)

	//se il file gobfile non esiste lo crea
	//gobfile è il file dove verrà resa persistente
	if _, err := os.Stat(gobfile); os.IsNotExist(err) {
		os.Create(gobfile)
	}

	//Carico dentro Contatori i dati salvati precedentemente
	err = load(gobfile, &F) //Carica in Contatori i dati salvati sul gobfile
	if err != nil {
		fmt.Println(err.Error())
	}

	//Per tutti i file passati come argomento esegue una goroutine
	for _, file := range os.Args[1:] {
		fmt.Println(file)
		wg.Add()
		go leggizip2(file)
	}

	wg.Wait() //Attende che terminino tutte le go routines

	/* var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err := e.Encode(Contatori); err != nil {
		panic(err)
	}
	fmt.Println("Encoded Struct ", b) */

	//Salva i dati in Contatori dentro il gobfile
	err = save(gobfile, F)
	if err != nil {
		log.Fatal(err.Error())
	}

	//Crea una variabile di tipo contatori per caricare tutti i dati salvati
	//nel gobfile
	var FruizioniDecoded Fruizioni
	err = load(gobfile, &FruizioniDecoded)
	if err != nil {
		log.Fatal(err.Error())
	}

	//fmt.Println(FruizioniDecoded)

	numFruizioni := len(FruizioniDecoded.Hashfruizione)
	fmt.Println(numFruizioni)
	for record := range FruizioniDecoded.Hashfruizione {

		fmt.Println(record)
		fmt.Println(FruizioniDecoded.Clientip[record])
		fmt.Println(FruizioniDecoded.Idvideoteca[record])
		speeds := FruizioniDecoded.Details[record]
		mean := stat.Mean(FruizioniDecoded.Details[record], nil)
		fmt.Printf("Media: %.3f\n", stat.Mean(speeds, nil))
		//harmonicmean := stat.HarmonicMean(speeds, nil)
		fmt.Printf("MediaArmonica: %.3f\n", stat.HarmonicMean(speeds, nil))
		mode, _ := stat.Mode(speeds, nil)
		fmt.Printf("Moda: %.3f\n", mode)
		nums := speeds
		sort.Float64s(nums) //Mette in ordine nums
		fmt.Printf("Mediana: %.3f\n", stat.Quantile(0.5, stat.Empirical, nums, nil))
		stdev := stat.StdDev(speeds, nil)
		fmt.Printf("StDev: %.3f\n", stat.StdDev(speeds, nil))
		fmt.Printf("Skew: %.3f\n", stat.Skew(speeds, nil))
		fmt.Printf("Curtosi: %.3f\n", stat.ExKurtosis(speeds, nil))

		fmt.Printf("NumChunks: %v\n", len(speeds))
		e := 0
		for _, n := range nums {
			var sigma float64
			sigma = 1 //tot sigma di distanza
			if (n-mean)/stdev < (-sigma * stdev) {
				e++
			}
		}
		if e > 0 {
			//se sono presenti errori ne mostra il quantitativo
			fmt.Printf("ERRORI: %d\n", e)
			fe := new(Fruizioniexport)
			fe.Hashfruizione = record
			fe.Clientip = FruizioniDecoded.Clientip[record]
			fe.Idvideoteca = FruizioniDecoded.Idvideoteca[record]
			fe.Errori = e

			l, err := json.Marshal(fe)
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Println(string(l))
			fmt.Println()
		}
		fmt.Println()

	}

	return
}
