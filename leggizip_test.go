package main

import (
	"testing"
)

// func ExampleLeggizip() {
// 	leggizip.Leggizip("we_ingestlog.gz")
// 	// Unordered output: 1
// 	// [18/Jun/2016:23:59:59.860+0000] http://vodos.oscdn.skycdn.it/RM/live/skycinema/MOVIE/skycinema-MVXY0000000000576770-201606141912150000.nff 185.26.141.212/ 185.26.141.212 41501 41501 979329565 100 5.515 1 206 video/nff No - 44909|HTTP|Sat_Jun_18_23:59:59_2016|0|1 SUCCESS_FINISH - - -
// }

func BenchmarkLeggizip(b *testing.B) {
	var wg = sizedwaitgroup.New(200)
	for n := 0; n < b.N; n++ {
		wg.Add()
		leggizip("we_ingestlog.gz")
	}
}
	}
}
