package main

/**
1. select + default 语句不会阻塞读取的操作，如果通道中没有数据，会立刻执行 default 中的代码。
*/
import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Start")

	message := make(chan string, 2)

	go (func() {
		message <- "hi"
	})()
	select {
	case msg := <-message:
		fmt.Printf("get message: %s \n", msg)
	case <-time.After(time.Second):
		fmt.Println("time out in 1s")
		// default:
		// 	fmt.Println("nothing receive")
	}

	fmt.Println("End")
}
