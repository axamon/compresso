package main

import (
	"flag"
	"sort"

	"github.com/Pallinder/go-randomdata"

	//"github.com/spf13/viper"

	"github.com/gonum/matrix/mat64"
	ma "github.com/mxmCherry/movavg"
	kalman "github.com/ryskiwt/go-kalman"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

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
			log.Panic(err.Error())
		}
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
			kalmansample(speeds, stdev)
			//fmt.Println()
		}
		//fmt.Println()

	}

	return
	//runtime.Goexit()

	//fmt.Println("Exit")
}

func kalmansample(speeds []float64, stdev float64) {

	sma3 := ma.ThreadSafe(ma.NewSMA(3)) //creo una moving average a 3
	sma7 := ma.ThreadSafe(ma.NewSMA(7)) //creo una moving average a 3

	//
	// kalman filter
	//

	//sstd := 0.000001
	sstd := stdev
	ostd := 0.1

	// trend model
	filter, err := kalman.New(&kalman.Config{
		F: mat64.NewDense(2, 2, []float64{2, -1, 1, 0}),
		G: mat64.NewDense(2, 1, []float64{1, 0}),
		Q: mat64.NewDense(1, 1, []float64{sstd}),
		H: mat64.NewDense(1, 2, []float64{1, 0}),
		R: mat64.NewDense(1, 1, []float64{ostd}),
	})
	if err != nil {
		panic(err)
	}

	n := len(speeds)
	s := mat64.NewDense(1, n, nil)
	x, dx := 0.0, 0.01
	xary := make([]float64, 0, n)
	yaryOrig := make([]float64, 0, n)
	ma3 := make([]float64, 0, n)
	ma7 := make([]float64, 0, n)

	//aryOrig := speeds
	for i := 0; i < n; i++ {
		//y := math.Sin(x) + 0.1*(rand.NormFloat64()-0.5)
		y := speeds[i]
		s.Set(0, i, y)
		x += dx

		xary = append(xary, x)
		yaryOrig = append(yaryOrig, y)
		ma3 = append(ma3, sma3.Add(y)) //aggiung alla media mobile il nuovo valore e storo la media
		ma7 = append(ma7, sma7.Add(y)) //aggiung alla media mobile il nuovo valore e storo la media
		//Verifica anomalia
		if ma3[i] < ma7[i] {
			fmt.Fprint(os.Stderr, "violazione soglia")
		}
	}

	filtered := filter.Filter(s)
	yaryFilt := mat64.Row(nil, 0, filtered)

	//
	// plot
	//

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	err = plotutil.AddLinePoints(p,
		"Original", generatePoints(xary, yaryOrig),
		"Filtered", generatePoints(xary, yaryFilt),
		"MA3", generatePoints(xary, ma3),
		"MA7", generatePoints(xary, ma7),
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	name := randomdata.FirstName(1) //da cambiare
	if err := p.Save(8*vg.Inch, 4*vg.Inch, name+".png"); err != nil {
		panic(err)
	}
}

func generatePoints(x []float64, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))

	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	return pts
}
