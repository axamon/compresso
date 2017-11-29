package main

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

// func ExampleLeggizip() {
// 	leggizip.Leggizip("we_ingestlog.gz")
// 	// Unordered output: 1
// 	// [18/Jun/2016:23:59:59.860+0000] http://vodos.oscdn.skycdn.it/RM/live/skycinema/MOVIE/skycinema-MVXY0000000000576770-201606141912150000.nff 185.26.141.212/ 185.26.141.212 41501 41501 979329565 100 5.515 1 206 video/nff No - 44909|HTTP|Sat_Jun_18_23:59:59_2016|0|1 SUCCESS_FINISH - - -
// }

func ExampleNewClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleLeggizip() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var Wg = sync.WaitGroup
	fmt.Println(Test)
	Wg.Add(1)
	go leggizip("we_ingestlog_clf_81.74.224.5_20160619_000000_52234.gz")
	Wg.Wait()
	val, err := client.SCard("recordhashes").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val, err)
	// Output: 96 <nil>
}
