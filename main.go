package main

import (
	"sort"

	"gonum.org/v1/gonum/stat"

	//"compress/gzip"

	"encoding/json"
	"fmt"
	"log"
	"os"

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
	Edgeip        string
	Giorno        string
	Orario        string
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
			/* 	fe.Giorno = Giorno
			fe.Errori = e */

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
