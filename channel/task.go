/**
 4 workers, 10 tasks
 output like:

 worker 1: startup
 worker 2: startup
 worker 1: task 1 start
 worker 3: startup
 worker 2: task 2 start
 worker 4: startup
 worker 3: task 3 start
 worker 3: task 4 start
 //.......
 worker %d: task 1 completed
 // ......
 worker %d: task 10 completed

 worker %d: shutdown
**/

package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	wg.Add(4)

	// create 10 tasks
	tasks := make(chan int, 10)

	// create 4 workers
	for i := 1; i <= 4; i++ {
		go worker(tasks, i)
	}
	for i := 1; i <= 10; i++ {
		tasks <- i
	}

	close(tasks)
	wg.Wait()
}

func worker(tasks chan int, i int) {
	fmt.Printf(" worker %d: startup\n", i)
	defer wg.Done()
	for {
		task, ok := <-tasks
		if !ok {
			// all done
			fmt.Printf(" worker %d: shutdown\n", i)
			return
		}
		fmt.Printf(" worker %d: task %d start\n", i, task)
		time.Sleep(2 * time.Second)
		fmt.Printf(" worker %d: task %d completed\n", i, task)
	}
}
