package main

import (
	"container/ring"
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"

	randomdata "github.com/Pallinder/go-randomdata"
)

func inserisci(num int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(in)
	for n := 1; n <= num; n++ {
		aggettivo := randomdata.Adjective()
		//fmt.Println(aggettivo) //debug
		in <- aggettivo
		//fmt.Println("Inserito")
	}
}

func elabora(wg *sync.WaitGroup) {
	defer wg.Done()
	//defer close(out)
	for i := range in {
		h := md5.Sum([]byte(i))
		ringlock.Lock()
		r.Value = h
		r = r.Next()
		ringlock.Unlock()
		//fmt.Println("Elaborato")
	}
}

var num int
var ringlock sync.Mutex
var r *ring.Ring

var in = make(chan string, num)

//var out = make(chan [16]byte, 1)

func main() {
	ctx, cancel := context.WithCancel(context.Background()) //crea un context globale
	defer cancel()

	go func() {
		// If ctrl+c is pressed it saves the situation and exit cleanly
		c := make(chan os.Signal, 1) //crea un canale con buffer unitario
		signal.Notify(c, os.Interrupt)
		s := <-c
		fmt.Println("Got signal:", s)
		ctx.Done()
		cancel() //fa terminare il context background
		//esce pulito
		fmt.Println("Uscita pulita")
		os.Exit(0)
	}()

	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err.Error())
	}

	r = ring.New(num)

	var wg sync.WaitGroup

	wg.Add(1)
	go inserisci(num, &wg)

	for i := 1; i <= num; i++ {
		wg.Add(1)
		go elabora(&wg)
	}

	wg.Wait()

	r.Do(
		func(p interface{}) {
			fmt.Printf("%x\n", p)
		})

	return

}
