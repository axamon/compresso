package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"runtime"
	"sync"
)

const gobfile = "./contatori.gob"

var gobfileLock sync.RWMutex

// Encode via Gob to file
func save(path string, object interface{}) error {
	gobfileLock.Lock()
	defer gobfileLock.Unlock()
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
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
