package main

import (
	"log"
	"time"
)

type Runner struct {
	tasks []func(int)
}

func New() *Runner {
	return &Runner{}
}

func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) Start() error {

	for i, task := range r.tasks {
		task(i)
	}

	return nil
}

func createTask() func(int) {
	return func(i int) {
		log.Printf("task %d executed\n", i)
		time.Sleep(2 * time.Second)
	}
}

func main() {
	log.Println("Start up")

	runner := New()

	runner.Add(createTask(), createTask(), createTask())

	if err := runner.Start(); err != nil {
		log.Fatal(err)
	}

}
