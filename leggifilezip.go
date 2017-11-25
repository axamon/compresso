package main

import (
	"compresso/leggizipevolved"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	leggizipevolved.LeggizipEvolved(os.Args[1])
}
