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

	done := make(chan bool)
	player := make(chan int, 2)
	// create the team
	go runner(player, done)
	i := 0
	for i < total {
		i++
		player <- i

	}

	close(player)
	<-done
	fmt.Println("done")

}

func runner(player chan int, done chan bool) {
	for {
		// get the baton
		value, ok := <-player
		if !ok {
			done <- true
			return
		}

		fmt.Printf("Runner %d running\n", value)

		time.Sleep(1 * time.Second)
	}
}
