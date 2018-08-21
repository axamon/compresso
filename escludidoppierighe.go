package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
)

var escludidoppionilock sync.RWMutex

var hashline map[string]bool

func escludidoppioni(line string) (err error) {
	err = nil
	hashriga := md5.New()                         //prepara a fare un hash
	hashriga.Write([]byte(line))                  //hasha tutta la linea
	Hash := hex.EncodeToString(hashriga.Sum(nil)) //estrae l'hash md5sum in versione quasi human readable
	escludidoppionilock.Lock()
	defer escludidoppionilock.Unlock()
	if hashline[Hash] == true {
		err = fmt.Errorf("Linea già inserita")
	}
	hashline[Hash] = true
	return err
}
