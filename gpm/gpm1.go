package main

import (
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)
	ch := make(chan bool)
	go func() {
		time.Sleep(1000000)
		ch <- true
	}()
	go func() {
		time.Sleep(1000)

	}()
	<-ch
}
