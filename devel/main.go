package main

import (
	"crypto/md5"
	"fmt"
	"sync"

	randomdata "github.com/Pallinder/go-randomdata"
)

var queueraw = make(chan string, 1)
var queuehashed = make(chan string, 1)
var queuepersisted = make(chan string, 1)

func randomentry() {
	for {
		queueraw <- randomdata.SillyName()
		queueraw <- randomdata.IpV4Address()
	}
}

func hasher() {
	for v := range queueraw {
		md5hash := md5.Sum([]byte(v))
		queuehashed <- fmt.Sprintf("%x", md5hash)
	}
}

var mappahash = make(map[string]bool)

var lock = sync.Mutex{}

func persist() {
	for v := range queuehashed {
		lock.Lock()
		mappahash[v] = true
		lock.Unlock()
	}

}

func main() {
	//	ctx, cancel := context.WithCancel(context.Background())
	//	defer cancel()

	go randomentry()
	go hasher()
	go persist()

	for v := range mappahash {
		fmt.Println(v)
	}

}
