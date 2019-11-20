package main

/**
panic 会导致『恐慌』，不同于Java中的 异常。
执行，且只执行，当前goroutine的defer。

panic有操守，退出前会执行本goroutine的defer，方式是原路返回(reverse order)
panic不多管，不是本goroutine的defer，不执行


*/
import "fmt"

func cleanup() {
	if r := recover(); r != nil {
		fmt.Println("cleanup:", r)
	}
}

func foo() {
	defer fmt.Println("foo start")
	fmt.Println("foo doing")
	panic("error foo")
	defer fmt.Println("foo end")
}

func main() {
	defer cleanup()
	fmt.Println("start")
	defer fmt.Println("defer before foo")
	foo()
	defer fmt.Println("defer after foo")
	fmt.Println("end")
}
