//  
/*
select 随机性解答： select 会随机选择一个可用通用做收发操作。 
所以代码是有肯触发异常，也有可能不会。 
单个 chan 如果无缓冲时，将会阻塞。
但结合 select 可以在多个 chan 间等待执行。
有三点原则： 
* select 中只要有一个 case 能 return，则立刻执行。 
* 当如果同一时间有多个 case 均能 return 则伪随机方式抽取任意一个执行。 
* 如果没有一个 case 能 return 则可以执行”default” 块。
*/
package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	intChan := make(chan int, 1)
	strChan := make(chan string, 1)
	intChan <- 1
	strChan <- "hello"
	select {
	case value := <-intChan:
		fmt.Println(value)
	case value := <-strChan:
		panic(value)
	}
}
