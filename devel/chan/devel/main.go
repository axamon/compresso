package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Pallinder/go-randomdata"
)

func inserisci(num int) chan string {
	out := make(chan string, num)
	//defer close(out)

	for n := 1; n <= num; n++ {
		aggettivo := randomdata.Adjective()
		//fmt.Println(aggettivo) //debug
		out <- aggettivo
	}
	return out
}

var num int

func main() {
	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err.Error())
	}

	c := inserisci(num)
	close(c)
	go func() {
		for v := range c {
			fmt.Printf("%v\n", v)
		}
	}()
	time.Sleep(3 * time.Second)
	fmt.Println("ok")
	for v := range c {
		fmt.Printf("%v\n", v)
	}

	//fmt.Println("ok")

	return

}
