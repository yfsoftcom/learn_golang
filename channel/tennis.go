/**
use channel to implement tennis

当前的版本有bug，会出现一个玩家自己和自己玩。
输出：
Hit round: 1
Player: Jack got it!
Hit round: 2
Player: Jack got it!
Hit round: 3
Player: Jack got it!
Hit round: 4
Player: Jack got it!
Hit round: 5
Player: Jack got it!
Hit round: 6
Player: Bob missed!
Player: Jack win!
因为2个goroutine之间没有互斥，可能其中一个会不断获得channel的数据。
修改方案，将 court := make(chan int, 1) 修改为 unbuffered channel 即可。
**/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

var done chan bool

func init() {
	rand.Seed(time.Now().UnixNano())
	done = make(chan bool)
}

func main() {

	// create the court for playground
	// 这里一定要使用 unbuffered chan
	court := make(chan int)

	// create 2 players
	go player("Bob", court)
	go player("Jack", court)

	// start
	court <- 1
	// wait game
	<-done
}

func player(name string, court chan int) {
	for {
		// 阻塞获取球
		counter, ok := <-court
		if !ok {
			// 对方 miss 了，win！
			fmt.Printf("Player: %s win!\n", name)
			done <- true
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
