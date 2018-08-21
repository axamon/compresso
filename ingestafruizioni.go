package main

import "sync"

var lock sync.RWMutex
var detailslock sync.RWMutex

func ingestafruizioni(hashfruizione, clientip, idvideoteca string, speed float64) {
	if F.Hashfruizione[hashfruizione] == false {
		lock.Lock()
		F.Hashfruizione[hashfruizione] = true
		F.Clientip[hashfruizione] = clientip
		F.Idvideoteca[hashfruizione] = idvideoteca
		lock.Unlock()
	}
	detailslock.Lock()
	F.Details[hashfruizione] = append(F.Details[hashfruizione], speed)
	detailslock.Unlock()
	return
}
