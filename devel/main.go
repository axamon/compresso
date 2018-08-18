package main

import (
	"fmt"
	"sync"
)

type Person struct {
	name string
	age  int
}

type People map[string]*Person

func main() {
	p := make(People)
	p["HM"] = &Person{"Hank McNamara", 39}
	p["HM"].age += 1
	fmt.Printf("age: %d\n", p["HM"].age)
	var counter = struct {
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}

}
