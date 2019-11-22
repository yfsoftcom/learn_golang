package main

import (
	"fmt"
	"sync"
	"time"
)

/**
使用 channel 和 goroutine 的组合，可以很好的实现多任务的并发协作
但是，对于一些共享变量，我们仍然需要进行一些数据的同步锁，防止多协程引发典型的 负库存 问题

sync.Mutex 可以帮助我们解决数据访问的问题，类似于 java 中的 Lock或者syncnized
**/

type Counter struct {
	val int
	mux sync.Mutex
}

func (c *Counter) safeIncr() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.val++
}

func (c *Counter) incr() {
	c.val++
}

func (c *Counter) reset() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.val = 0
}

func main() {
	fmt.Println("start")
	c := &Counter{val: 0}
	// do some incr or decr times
	for i := 0; i < 1000; i++ {
		go c.incr() // 这种方式下，得到的综合
	}

	time.Sleep(time.Second)
	fmt.Println(c.val)
	c.reset()
	for i := 0; i < 1000; i++ {
		go c.safeIncr()
	}

	time.Sleep(time.Second)
	fmt.Println(c.val)
	fmt.Println("end")
}
