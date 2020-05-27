/*

1. 使用channel进行协程的控制
2. for {
		msg, open := <-c
		if !open {
			break
		}
		fmt.Println(msg)
	}
	等同于
	for msg := range c {
		fmt.Println(msg)
	}
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("foo bar")
	c := make(chan string)
	go count("bob", c)
	go count("jack", c)

	for msg := range c {
		fmt.Println(msg)
	}

	fmt.Println("finish")
}

func count(name string, c chan string) {
	defer close(c)
	for i := 1; i <= 5; i++ {
		c <- name
		time.Sleep(time.Second * 1)
	}

}
