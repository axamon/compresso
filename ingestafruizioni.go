package main

import "sync"

var lock sync.RWMutex

//Fruizione conserva i dati utente delle singole fruizioni
type Fruizione struct { //creo un type dove mettere i dati di ogni singola fruizione
	hashfruizione string
	clientip      string
	idvideoteca   string
}

func ingestafruizioni(hashfruizione string, speed float64) {
	lock.Lock()
	Contatori.Fruizioni[hashfruizione] = true
	Contatori.Numchunks[hashfruizione]++
	Contatori.Details[hashfruizione] = append(Contatori.Details[hashfruizione], speed)
	lock.Unlock()
	return
}
