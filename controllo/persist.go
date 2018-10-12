package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"runtime"
	"sync"
)

const gobfile = "./gob/contatori.gob"

var gobfileLock = sync.RWMutex{}

// Encode via Gob to file
func save(path string, object interface{}) error {
	file, err := os.Create(path)
	defer file.Close()
	if err == nil {
		encoder := gob.NewEncoder(file)
		gobfileLock.Lock()
		encoder.Encode(object)
		gobfileLock.Unlock()
	}

	return err
}

// Decode Gob file
func load(path string, object interface{}) error {
	gobfileLock.RLock()
	defer gobfileLock.RUnlock()
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func check(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		//os.Exit(1)
	}
}
