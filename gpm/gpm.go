package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(1)
	fmt.Println("start")

	wg := sync.WaitGroup{}

	total := 257 // 这里的数字在 257 的左右会产生不同的规律， 256 是本地队列的数量上限。超过则会在本地和全局队列中伪随机获取groutine来执行，导致乱序

	wg.Add(total)
	for i := 0; i < total; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println(">. i =", i)
		}(i)
	}
	wg.Wait()
	fmt.Println("End")
}
