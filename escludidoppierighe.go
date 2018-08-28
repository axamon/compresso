package main

import (
	"fmt"
	"sync"
)

const hashlinefile = "./gob/hashline.gob"

var escludidoppionilock sync.RWMutex

var hashline map[string]bool

func escludidoppioni(line string) (err error) {
	err = nil
	Hash := md5sumOfString(line)

	escludidoppionilock.Lock()
	if hashline[Hash] == true {
		err = fmt.Errorf("Linea già inserita")
	}
	hashline[Hash] = true
	escludidoppionilock.Unlock()

	return err
}
