/*

1. 使用 waitgroup 的方法来实现并发的阻塞，但是需要使用 WaitGroup 来做并发计数，对于较小的项目适用
2. 需要在 go 起协程的时候使用匿名函数来进行 WaitGroup.Done() ，从整体上来说，不适合进行多模块的开发
3. WaitGroup 计数器无法进行实例化，最好不要在函数之间传递该变量来进行 Done() & Wait() ,会造成不可控的局面
4. 推荐使用 channel 来实现协程间的协作
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("foo bar")

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		count("bob")
	}()
	go func() {
		defer wg.Done()
		count("jack")
	}()

	wg.Wait()
	// fmt.Scanln()
}

func count(name string) {
	for i := 1; i <= 5; i++ {
		fmt.Println(i, name)
		time.Sleep(time.Second * 1)
	}

}
