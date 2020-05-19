/**
模拟接力赛的demo程序
target:
///////
Runner 1 running with baton
Runner 2 ready
Runner 1 exchange with Runner 2
Runner 2 running with baton
Runner 3 ready
Runner 2 exchange with Runner 3
Runner 3 running with baton
Runner 4 ready
Runner 3 exchange with Runner 4
Runner 4 running with baton
Runner 4 finish the race
**/

package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var total = 4

func main() {
	// create the baton
	baton := make(chan int, 1)

	wg.Add(1)
	// create the team
	go runner(baton)

	// start with 1
	baton <- 1
	wg.Wait()
}

func runner(baton chan int) {

	// get the baton
	value, ok := <-baton
	if !ok {
		return
	}

	fmt.Printf("Runner %d running with baton\n", value)

	time.Sleep(1 * time.Second)
	if value == total {
		// finish
		fmt.Printf("Runner %d finish the race\n", value)
		close(baton)
		wg.Done()
		return
	}
	newRunner := value + 1
	fmt.Printf("Runner %d ready\n", newRunner)
	go runner(baton)

	fmt.Printf("Runner %d exchange with Runner %d\n", value, newRunner)
	baton <- newRunner
}
