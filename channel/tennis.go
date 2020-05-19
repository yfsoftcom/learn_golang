/**
use channel to implement tennis
**/
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	// create the court for playground
	court := make(chan int, 1)

	// create 2 players
	wg.Add(2)
	go player("Bob", court)
	go player("Jack", court)

	// start
	court <- 1
	wg.Wait()
}

func player(name string, court chan int) {
	defer wg.Done()
	for {
		// 阻塞获取球
		counter, ok := <-court
		if !ok {
			// 对方 miss 了，win！
			fmt.Printf("Player: %s win!\n", name)
			return
		}

		fmt.Printf("Hit round: %d \n", counter)
		n := rand.Intn(100)
		if n%13 == 0 {
			//如果出现 13 倍数，则mis
			fmt.Printf("Player: %s missed!\n", name)
			close(court)
			return
		}

		fmt.Printf("Player: %s got it!\n", name)
		court <- counter + 1
	}
}
