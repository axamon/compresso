package main

import "sync"

var lock sync.RWMutex

//var detailslock sync.RWMutex

func ingestafruizioni(hashfruizione, clientip, idvideoteca, idaps, edgeip, giorno, orario string, speed float64) {
	lock.Lock()
	defer lock.Unlock()
	if F.Hashfruizione[hashfruizione] == false {
		F.Hashfruizione[hashfruizione] = true
		F.Clientip[hashfruizione] = clientip
		F.Idvideoteca[hashfruizione] = idvideoteca
		F.Idaps[hashfruizione] = idaps
		F.Edgeip[hashfruizione] = edgeip
		F.Giorno[hashfruizione] = giorno
		F.Orario[hashfruizione] = orario
	}

	//detailslock.Lock()
	F.Details[hashfruizione] = append(F.Details[hashfruizione], speed)
	//detailslock.Unlock()

	return
}
