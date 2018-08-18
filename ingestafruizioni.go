package main

type Fruizione struct { //creo un type dove mettere i dati di ogni singola fruizione
	hashfruizione string
	clientip      string
	idvideoteca   string
}

func ingestafruizioni(hashfruizione string, speed float64) {
	Contatori.Lock()
	Contatori.fruizioni[hashfruizione] = true
	Contatori.numchunks[hashfruizione]++
	Contatori.details[hashfruizione] = append(Contatori.details[hashfruizione], speed)
	Contatori.Unlock()
	return
}
