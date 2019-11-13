/*

1. 使用channel进行协程的控制
2. channel的读取是会造成阻塞的
3. 使用2个协程，用不同的频率向2个chan 中推送数据，然后在 main 中读取2个chan中接收到的数据，查看读取是否会造成写入的阻塞
4. ----
	a of 500ms
	b of 2000ms
	----
	程序先输出 a 但是必须等到 b 输出之后才会继续 输出 a ，以证明，通道的读取会阻塞写入的操作，并且读的操作本身是阻塞的
5. 使用 select 来异步读取不同 channel ，不会造成阻塞
	输出结果如下：
	b of 2000ms
	----
	a of 500ms
	----
	a of 500ms
	----
	a of 500ms
	----
	a of 500ms
	----
	b of 2000ms
	----
	Channel B 的读取不会造成 A 的阻塞，达到了完全的异步
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("foo bar")
	c1 := make(chan string)
	c2 := make(chan string)

	// go
	go func() {
		for {
			c1 <- "a of 500ms"
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			c2 <- "b of 2000ms"
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		select {
		case msg1 := <-c1:
			fmt.Println(msg1)
		case msg2 := <-c2:
			fmt.Println(msg2)
		}
		fmt.Println("----")
	}
	fmt.Println("finish")
}
