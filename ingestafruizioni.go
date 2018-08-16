package main

import "sync"

type fruizioni struct { //creo un type dove mettere i dati di ogni singola fruizione
	clientipint    int
	idvideoteca    int
	sumspeeds      int
	sumquadspeeds  int
	numvideochunks int
}

var h = make(map[string]*fruizioni) //creo una mappa e la istanzio per dare un nome (hashstruizione) a oghi raccolta dati

var counter struct {
	sync.RWMutex                       //locka e slocka l'accesso ai dati
	h            map[string]*fruizioni //bisogna refenziale con * che indica una istanza e non il type stesso
}

func ingestafruizioni(hashfruizione string, clientipint, idvideoteca, speed int) {
	counter.Lock()
	counter.h[hashfruizione].idvideoteca = idvideoteca
	counter.h[hashfruizione].clientipint = clientipint
	counter.h[hashfruizione].numvideochunks++
	counter.h[hashfruizione].sumspeeds += speed
	counter.h[hashfruizione].sumquadspeeds += (speed * speed)
	counter.Unlock()

}
