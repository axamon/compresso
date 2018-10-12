package main

import (
	"flag"
	"sort"

	//"github.com/spf13/viper"

	"gonum.org/v1/gonum/stat"

	//"compress/gzip"

	"encoding/json"
	"fmt"
	"log"
	"os"
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

//F contiene tutte le informazioni delle varie fruizioni
//è la variabile che verrà resa persistente su disco
var F = Fruizioni{}

//var sigma float64

var sigma = flag.Float64("sigma", 3, "in numero di stdev da tenere in considerazione")

func main() {

	/* 	v := viper.New()
	   	v.SetConfigFile("compresso.yaml")
	   	err := v.ReadInConfig()
	   	if err != nil {
	   		log.Fatal(err.Error())
	   	}
	   	//sigma := v.GetFloat64("sigma") */
	flag.Parse()

	/* F.Hashfruizione = make(map[string]bool)
	F.Clientip = make(map[string]string)
	F.Idvideoteca = make(map[string]string)
	F.Edgeip = make(map[string]string)
	F.Giorno = make(map[string]string)
	F.Orario = make(map[string]string)
	F.Details = make(map[string][]float64) */

	//se il file gobfile non esiste lo crea
	//gobfile è il file dove verrà resa persistente
	if _, err := os.Stat(gobfile); os.IsNotExist(err) {
		os.Create(gobfile)
	}

	//Crea una variabile di tipo contatori per caricare tutti i dati salvati
	//nel gobfile
	var FruizioniDecoded Fruizioni
	err := load(gobfile, &FruizioniDecoded)
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
		var numchunks int
		numchunks = len(speeds)

		mean := stat.Mean(FruizioniDecoded.Details[record], nil)
		//fmt.Printf("Media: %.3f\n", stat.Mean(speeds, nil))
		harmonicmean := stat.HarmonicMean(speeds, nil)

		//fmt.Printf("MediaArmonica: %.3f\n", stat.HarmonicMean(speeds, nil))
		mode, _ := stat.Mode(speeds, nil)
		//mean, _ := stat.Mode(speeds, nil)

		//fmt.Printf("Moda: %.3f\n", mode)
		nums := speeds
		//fmt.Println(len(nums))
		entropy := stat.Entropy(nums)
		sort.Float64s(nums) //Mette in ordine nums
		//fmt.Printf("Mediana: %.3f\n", stat.Quantile(0.5, stat.Empirical, nums, nil))
		median := stat.Quantile(0.5, stat.Empirical, nums, nil)
		percentile95 := stat.Quantile(0.95, stat.Empirical, nums, nil)

		stdev := stat.StdDev(speeds, nil)
		stderr := stat.StdErr(stdev, float64(numchunks))
		//fmt.Printf("StDev: %.3f\n", stat.StdDev(speeds, nil))
		skew := stat.Skew(speeds, nil)
		//fmt.Printf("Skew: %.3f\n", stat.Skew(speeds, nil))
		curtosi := stat.ExKurtosis(speeds, nil)
		chisquare := stat.ChiSquare(nums, speeds)
		//fmt.Printf("Curtosi: %.3f\n", stat.ExKurtosis(speeds, nil))

		//fmt.Printf("NumChunks: %v\n", len(speeds))
		//e := -1 //un errore lo abboniamo
		e := 0
		s := *sigma
		lowerlimit := -s * stdev
		for _, n := range nums {
			x := n - harmonicmean

			if x < lowerlimit {
				e++

			} else {
				break
			}
		}
		if e > 0 {
			//se sono presenti errori ne mostra il quantitativo
			//fmt.Printf("ERRORI: %d\n", e)
			fe := new(Fruizioniexport)
			fe.Hashfruizione = record
			fe.Clientip = FruizioniDecoded.Clientip[record]
			fe.Idvideoteca = FruizioniDecoded.Idvideoteca[record]
			fe.Idaps = FruizioniDecoded.Idaps[record]
			fe.Edgeip = FruizioniDecoded.Edgeip[record]
			fe.Giorno = FruizioniDecoded.Giorno[record]
			fe.Orario = FruizioniDecoded.Orario[record]
			fe.Media = mean
			fe.Stdev = stdev
			fe.Moda = mode
			fe.Mediana = median
			fe.MediaArmonica = harmonicmean
			fe.Percentile95 = percentile95
			fe.Skew = skew
			fe.Curtosi = curtosi
			fe.Numchunks = numchunks
			fe.Stderr = stderr
			fe.Errori = e
			fe.Entropia = entropy
			fe.Chisquare = chisquare

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
