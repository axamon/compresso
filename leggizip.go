package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"log"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`(\bcb.ticdn.it\b|\bvodm0.ticdn.it\b)`)

func leggizip2(ctx context.Context, file string) {
	defer wg.Done()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1) //esegue una go routine su tutti i processori -1

	err := escludidoppioni(ctx, file)
	if err != nil {
		log.Printf("file %s già elaborato", file)
		return
	}

	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		log.Println(err.Error())
	}

	gr, err := gzip.NewReader(f)
	//gr, err := pgzip.NewReaderN(f, 4096, 100) //sfrutta il gzip con steroide che legge nel futuro per andare più veloce assai
	defer gr.Close()
	if err != nil { //se però si impippa qualcosa allora blocca tutto
		log.Println(err.Error())
		return
	}

	fileelements := strings.Split(file, "_") //prende il nome del file di log e recupera i campi utili
	Type := fileelements[1]                  //qui prede il tipo di log
	edgeip := fileelements[3]
	giorno := fileelements[4]
	orario := fileelements[5] //qui prende l'ip della cache

	if Type == "accesslog" { //se il tipo di log è "accesslog" allora fa qualcosa che ancora non ho finito di fare

		log.Printf("Elaboro log %s", file)
		defer log.Printf("Finita elaborazione log %s", file)

		scan := bufio.NewScanner(gr)
		for scan.Scan() {
			line := scan.Text()
			/* err := escludidoppioni(ctx, line) //tocca toglierlo che crea file troppo grossi co tutti sti hash
			if err != nil {
				continue
			} */
			if !strings.HasPrefix(line, "[") { //se la linea non inzia con [ allora salta
				continue
			}

			s := strings.Split(line, "\t")

			if len(s) < 8 || !re.MatchString(s[6]) { //se s[6] il campo url non contiene quanto compilato in re ignora la riga
				continue
			}

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
				log.Println(err.Error())
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
			//status := s[3] //da usare per errori 40x e 50x
			ua := s[8]

			//fmt.Println(Urlschema)
			//Urlhost := u.Host
			Urlpath := u.Path
			//fmt.Println(Urlpath)
			//Urlquery := u.RawQuery
			//Urlfragment := u.Fragment
			pezziurl := strings.Split(Urlpath, "/")
			//fmt.Println(pezziurl)
			if len(pezziurl) < 11 {
				continue
			}
			if ok := strings.Contains(Urlpath, "video="); ok == true { //solo i chunk video
				//fmt.Println(pezziurl)
				//fmt.Println(Urlpath)
				idvideoteca := pezziurl[6]
				//tipocodifica := pezziurl[7]
				idaps := pezziurl[8]
				//fmt.Println(idvideoteca)
				//encoding := pezziurl[10]
				//fmt.Println(encoding)
				//re := regexp.MustCompile(`QualityLevels\(([0-9]+)\)$`)
				//bitratestr := re.FindStringSubmatch(encoding)[1]
				//bitrate, _ := strconv.ParseFloat(bitratestr, 8)
				/* if err != nil {
					log.Fatal(err.Error())
				} */
				//bitrateMB := bitrate * bitstoMB
				Hash := md5sumOfString(ctx, clientip+idvideoteca+ua)

				ingestafruizioni(Hash, clientip, idvideoteca, idaps, edgeip, giorno, orario, speed)
			}
			if ok := strings.Contains(Urlpath, "DASH"); ok == true { //Prende solo i chunk DASH
				idvideoteca := pezziurl[6]
				idaps := pezziurl[8]
				Hash := md5sumOfString(ctx, clientip+idvideoteca+ua)

				ingestafruizioni(Hash, clientip, idvideoteca, idaps, edgeip, giorno, orario, speed)
			}
		}
	}
	return //terminata la Go routine!!! :)
}
