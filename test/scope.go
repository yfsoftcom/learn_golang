// test the scope
//
/**
start
2. i = 9
1. i = 10
1. i = 10
1. i = 10
1. i = 10
1. i = 10
1. i = 10
1. i = 10
1. i = 10
1. i = 10
1. i = 10
2. i = 0
2. i = 1
2. i = 2
2. i = 3
2. i = 4
2. i = 5
2. i = 6
2. i = 7
2. i = 8
End
*/
//
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

	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			fmt.Println("1. i =", i)
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println("2. i =", i)
		}(i)
	}
	wg.Wait()
	fmt.Println("End")
}
