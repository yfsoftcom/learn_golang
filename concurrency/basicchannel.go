/*

Dont communicate by sharing memory, instead, share memory by communicating.
不用通过共享内存来通信，我们应该通过通信来共享内存。

By default, sends and receives block until the other side is ready.
This allows goroutines to synchronize without explicit locks or condition variables.

通常，数据在发送到通道中之后就丧失了该数据的控制权，也只会被一个 goroutine 接受到该通道中的数据，所以可以避免协程间资源竞争的问题。
但是，如果通道中的数据是 指针的话，发送方发送之后同样保存了该值的引用，同样可以访问该资源，引起资源的竞争。
所以，不正确的使用 channel 同样会导致漏洞。

1. 使用channel进行数据的存取
2. 不使用 goroutine 的情况和 使用 goroutine 的情况？
		无缓冲的 channel 是无法在 非 goroutine 的情况下使用的，也就是说一定要在 go 启动的非主线程中的操作才行。
		带缓冲的 channel 可以。

3. 由于 goroutine 的执行时机是随机的，所以使用 channel 获取通道数据的时候，取到的数据顺序也是随机的。
4. 使用 select case default 读取channel
5. 通道做为参数是可以定义方向的，定义该参数只能够发送数据还是只能够接收数据。
6. we can get the default value after the channel closed.
   ```
   c := make(chan int)
   go () {
	   c <- 100
   }
   close(c)
   v := <-c 	// v is 0(the channel's default value)
   ```
参考引用

https://golang.org/doc/effective_go.html#channels
https://draveness.me/whys-the-design-communication-shared-memory
https://gobyexample.com/channel-directions
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("foo bar")
	c := make(chan int, 1)

	go foo(c, 1)
	c <- 50
	c <- 500
	close(c)
	// for i := 1; i <= 10; i++ {
	// 	go foo(c, i)
	// }
	// time.Sleep(time.Second * 2)
	v := <-c // 这里会发生阻塞
	fmt.Println(v)
	// close(c)
	fmt.Println("finish")
}

/**
accept c channel int,
**/
func foo(c chan int, i int) {
	time.Sleep(time.Second * 2)
	c <- i * 100
}
