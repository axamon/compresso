package main

type Fruizione struct { //creo un type dove mettere i dati di ogni singola fruizione
	hashfruizione string
	clientip      string
	idvideoteca   string
}

//var details = make(map[string][]float64)
//var fruizioni = make(map[string]bool)

//var numchunks = make(map[string]int)

//var sumspeeds = make(map[string]float64)
//var sumsquarespeeds = make(map[string]float64)

func ingestafruizioni(hashfruizione string, speed float64) {
	Contatori.Lock()
	Contatori.fruizioni[hashfruizione] = true
	Contatori.numchunks[hashfruizione]++
	Contatori.details[hashfruizione] = append(Contatori.details[hashfruizione], speed)
	//sumspeeds[hashfruizione] += speed
	//sumsquarespeeds[hashfruizione] += (speed * speed)
	Contatori.Unlock()
	return
}
