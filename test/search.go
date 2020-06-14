// 已知一个足够大的切片，类型为 int
// 通过多个goroutine进行搜索，找到给定的一个 数据
// 要求：不可以进行排序，goroutine数量必须大于1
// 找到后，输出 Found it，并取消所有其他的 goroutine
// 超过 给定的 时间之后，输出 timeout，并关闭所有的 goroutine
package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

const (
	total   = 10000
	workers = 5
	max     = 100
	target  = 50
	timeout = 2 * time.Second
)

var result chan bool

func genData() *[]int {
	data := make([]int, total)
	for i := 0; i < total; i++ {
		data[i] = rand.Intn(max)
	}
	return &data
}

func search(ctx context.Context, result chan bool, data *[]int, start int, end int) {

	_, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := start; i < end; i++ {
		d := (*data)[i]
		fmt.Println(d, target)
		if target == d {
			result <- true
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

}

func main() {
	fmt.Println("start")

	start := time.Now()
	defer (func() {
		dur := time.Since(start)
		fmt.Println("total:", dur)
	})()
	numbPtr := flag.Int("numb", target, "an int")
	result = make(chan bool, 1)

	// init the slice
	data := genData()
	(*data)[rand.Intn(total)] = *numbPtr
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rows := total / workers
	for i := 0; i < workers; i++ {
		go search(ctx, result, data, i*rows, (i+1)*rows)
	}

	select {
	case <-result:
		fmt.Println("Found it")
	case <-ctx.Done():
		fmt.Println("Timeout", ctx.Err())
	}
	close(result)
}
