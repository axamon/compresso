package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/klauspost/pgzip"
)

func leggizip2(file string) {
	defer wg.Done()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1) //esegue una go routine su tutti i processori -1

	/* 	client := redis.NewClient(&redis.Options{ //connettiti a Redis server
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}) */

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
	/* SEIp := fileelements[3]
	Giorno := fileelements[4]
	Orario := fileelements[5] //qui prende l'ip della cache */

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
				Hash := md5sumOfString(clientip + idvideoteca + ua)
				/* hasher := md5.New()                                 //prepara a fare un hash
				hasher.Write([]byte(clientip + idvideoteca + ua))   //hasha tutta la linea
				Hash := hex.EncodeToString(hasher.Sum(nil))         //estrae l'hash md5sum in versione quasi human readable
				_, err = client.SAdd("recordhashes", Hash).Result() //finalmente usiamo l'hash dentro a redis
				if err != nil {
					log.Fatal(err.Error())
				} */
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
				Hash := md5sumOfString(clientip + idvideoteca + ua)
				/* 	hasher := md5.New()                                 //prepara a fare un hash
				hasher.Write([]byte())                              //hasha tutta la linea
				Hash := hex.EncodeToString(hasher.Sum(nil))         //estrae l'hash md5sum in versione quasi human readable
				_, err = client.SAdd("recordhashes", Hash).Result() //finalmente usiamo l'hash dentro a redis
				if err != nil {
					log.Fatal(err.Error())
				} */
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
