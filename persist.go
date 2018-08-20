package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"runtime"
	"sync"
)

const gobfile = "./contatori.gob"

var GobfileLock sync.RWMutex

// Encode via Gob to file
func Save(path string, object interface{}) error {
	GobfileLock.Lock()
	defer GobfileLock.Unlock()
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	GobfileLock.RLock()
	defer GobfileLock.RUnlock()
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func Check(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		//os.Exit(1)
	}
}
