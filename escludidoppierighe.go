package main

import (
	"context"
	"fmt"
	"sync"
)

const hashlinefile = "./gob/hashline.gob"

var escludidoppionilock = sync.RWMutex{}

var hashline map[string]bool

func escludidoppioni(ctx context.Context, line string) (err error) {
	err = nil
	Hash := md5sumOfString(ctx, line)

	escludidoppionilock.Lock()
	if hashline[Hash] == true {
		err = fmt.Errorf("Linea già inserita")
	}
	hashline[Hash] = true
	escludidoppionilock.Unlock()
	return err
}
