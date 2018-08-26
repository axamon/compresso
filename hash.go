package main

import (
	"crypto/md5"
	"fmt"
)

func md5sumOfString(value string) string {

	md5hash := md5.Sum([]byte(value))
	return fmt.Sprintf("%x", md5hash)

}
