/**
runner 任务执行工具的使用
**/
package main

import (
	runner "advanced-test/runner"
	"log"
	"time"
)

func main() {
	log.Println("start!")

	timeCost := 3 * time.Second

	scheRunner := runner.New(timeCost)

	scheRunner.Add(createTask(), createTask(), createTask())

	scheRunner.Start()

}

func createTask() func(int) {
	return func(id int) {
		log.Printf("task %d running\n", id)
	}
}
