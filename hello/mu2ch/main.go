package main

//
// Mutex与Channel性能比较
//
// Channel的性能与缓存有关。
// 0 缓存（~3）
// 	Mutex: 16777215 2.4276197s
// 	Channel: 16777215 7.1257563s
// 1 缓存（~2）
// 	Mutex: 16777215 2.3945976s
// 	Channel: 16777215 4.5510383s
// 10 缓存（~1.5）
// 	Mutex: 16777215 2.6647792s
// 	Channel: 16777215 3.4272882s
// 30 缓存（~）
//	Mutex: 16777215 2.8328827s
//	Channel: 16777215 2.8799139s
// 50 缓存（略优）
// 	Channel: 16777215 2.8188805s
// 	Mutex: 16777215 2.8469001s

// 结论：
// - 如果缓存足够，可与Mutex媲美。但无需更大的缓存。
// - Channel的优势应该是逻辑，同步Channel没有性能优势。
//

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex
var ch = make(chan int, 40)
var max = 1 << 24

var muVal, chVal int

func muWrite(n int) {
	mu.Lock()
	muVal = n
	mu.Unlock()
}

func muRead() int {
	mu.Lock()
	defer mu.Unlock()
	return muVal
}

func chServ() {
	for i := 0; i < max; i++ {
		ch <- i
	}
	close(ch)
}

func chTest() {
	tm := time.Now()

	for v := range ch {
		chVal = v
	}
	fmt.Println("Channel:", chVal, time.Since(tm))
}

func muServ() {
	for i := 0; i < max; i++ {
		muWrite(i)
	}
}

func muTest() {
	tm := time.Now()

	for i := 0; i < max; i++ {
		muVal = muRead()
	}
	fmt.Println("Mutex:", muVal, time.Since(tm))
}

func main() {
	go chServ()
	go muServ()

	go chTest()
	go muTest()

	time.Sleep(10 * time.Second)
	fmt.Println("Done!")
}
