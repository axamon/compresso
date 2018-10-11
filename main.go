package main

import (
	"context"
	"os/signal"
	"sort"

	"github.com/spf13/viper"

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

func init() {
	dir := "./gob"                                  //directory di storage dei file gob
	if _, err := os.Stat(dir); os.IsNotExist(err) { //se la directory non esiste la crea
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat("compresso.yaml"); os.IsNotExist(err) { //se il file di configurazione non esiste lo crea
		f, err := os.Create("compresso.yaml")
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString("sigma: \"3\"") //insersce il parametro sigma a 3 come default
		if err != nil {
			panic(err)
		}
		f.Close()
	}

}

var wg = sizedwaitgroup.New(200) //massimo numero di go routine per volta

//F contiene tutte le informazioni delle varie fruizioni
//è la variabile che verrà resa persistente su disco
var F = Fruizioni{}

var sigma float64

func main() {
	ctx, cancel := context.WithCancel(context.Background()) //crea un context globale
	defer cancel()

	//hashline tiene l'hash di ogni singolo linea di log
	hashline = make(map[string]bool)

	// If ctrl+c is pressed it saves the situation and exit cleanly
	c := make(chan os.Signal, 1) //crea un canale con buffer unitario
	signal.Notify(c, os.Interrupt)

	//Goroutine per la gestione dell'uscita dal programma tramite ctrl+c
	go func() {
		s := <-c
		fmt.Println("Got signal:", s)
		ctx.Done()
		//wg.Wait() //Attende che terminino tutte le go routines
		cancel() //fa terminare il context background
		//salva le mappe come file .gob
		err := save(hashlinefile, hashline)
		if err != nil {
			fmt.Printf("Uscita con errori: %s\n", err.Error())
			os.Exit(1)
		}
		//esce pulito
		fmt.Println("Uscita pulita")
		os.Exit(0)
	}()

	v := viper.New()
	v.SetConfigFile("compresso.yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	sigma := v.GetFloat64("sigma")
	//flag.Parse()

	//se il file hashlinefile non esiste lo crea
	//gobfile è il file dove verrà resa persistente
	if _, err := os.Stat(hashlinefile); os.IsNotExist(err) {
		h, err := os.Create(hashlinefile)
		h.Close()
		if err != nil {
			log.Printf("impossibile creare %s\n", hashlinefile)
		}
	}

	//Carico dentro hashline i dati salvati precedentemente
	err = load(hashlinefile, &hashline) //Carica in Contatori i dati salvati sul gobfile
	if err != nil {
		fmt.Println(err.Error()) //se da errore forse manca il file... non importa se è il primo avvio
	}
	defer save(hashlinefile, hashline)

	//go watch(ctx)
	//time.Sleep(10 * time.Second)

	F.Hashfruizione = make(map[string]bool)
	F.Clientip = make(map[string]string)
	F.Idvideoteca = make(map[string]string)
	F.Edgeip = make(map[string]string)
	F.Giorno = make(map[string]string)
	F.Orario = make(map[string]string)
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
		//fmt.Println(file)
		wg.Add()
		go leggizip2(ctx, file)
		//	file := "examplelogs/we_accesslog_clf_81.74.227.47_20181006_033600_12216.gz"
		//leggizip2(ctx, file)
	}

	wg.Wait() //Attende che terminino tutte le go routines

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
	fmt.Printf("Verificate %v fruizioni\n", numFruizioni)
	for record := range FruizioniDecoded.Hashfruizione {

		//fmt.Println(record)
		//fmt.Println(FruizioniDecoded.Clientip[record])
		//fmt.Println(FruizioniDecoded.Idvideoteca[record])
		speeds := FruizioniDecoded.Details[record]
		mean := stat.Mean(FruizioniDecoded.Details[record], nil)
		//fmt.Printf("Media: %.3f\n", stat.Mean(speeds, nil))
		//harmonicmean := stat.HarmonicMean(speeds, nil)
		//fmt.Printf("MediaArmonica: %.3f\n", stat.HarmonicMean(speeds, nil))
		//mode, _ := stat.Mode(speeds, nil)
		//fmt.Printf("Moda: %.3f\n", mode)
		nums := speeds
		//fmt.Println(len(nums))
		sort.Float64s(nums) //Mette in ordine nums
		//fmt.Printf("Mediana: %.3f\n", stat.Quantile(0.5, stat.Empirical, nums, nil))
		stdev := stat.StdDev(speeds, nil)
		//fmt.Printf("StDev: %.3f\n", stat.StdDev(speeds, nil))
		//fmt.Printf("Skew: %.3f\n", stat.Skew(speeds, nil))
		//fmt.Printf("Curtosi: %.3f\n", stat.ExKurtosis(speeds, nil))

		//fmt.Printf("NumChunks: %v\n", len(speeds))
		e := 0
		for _, n := range nums {
			if (n-mean)/stdev < (-sigma * stdev) {
				e++
			}
		}
		if e > 0 {
			//se sono presenti errori ne mostra il quantitativo
			//fmt.Printf("ERRORI: %d\n", e)
			fe := new(Fruizioniexport)
			fe.Hashfruizione = record
			fe.Clientip = FruizioniDecoded.Clientip[record]
			fe.Idvideoteca = FruizioniDecoded.Idvideoteca[record]
			fe.Edgeip = FruizioniDecoded.Edgeip[record]
			fe.Giorno = FruizioniDecoded.Giorno[record]
			fe.Orario = FruizioniDecoded.Orario[record]
			fe.Errori = e

			l, err := json.Marshal(fe)
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Println(string(l))
			//fmt.Println()
		}
		//fmt.Println()

	}

	return
	//runtime.Goexit()

	//fmt.Println("Exit")
}
