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
	"time"
)

var total = 4

func main() {
	// create the baton
	baton := make(chan int, 1)
	doneSignal := make(chan bool)
	// create the team
	go runner(baton, doneSignal)

	// start with 1
	baton <- 1
	<-doneSignal
}

func runner(baton chan int, done chan bool) {

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
		done <- true
		return
	}
	newRunner := value + 1
	fmt.Printf("Runner %d ready\n", newRunner)
	go runner(baton, done)

	fmt.Printf("Runner %d exchange with Runner %d\n", value, newRunner)
	baton <- newRunner
}
