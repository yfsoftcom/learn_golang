package main

import "fmt"

func test1() (x *int32) {
	var n int32
	x = &n
	defer fmt.Printf("in test case x1=%d\n", *x)
	*x = 7
	return
}

func test2() (x int) {
	x = 7
	defer fmt.Printf("in test case x2=%d\n", x)
	return 9
}

func test3() (x int) {
	defer func() {
		fmt.Printf("in test case x3=%d\n", x)
	}()
	x = 7
	return 9
}

func test4() (x int) {
	defer func(n int) {
		fmt.Printf("in test case x4_in=%d\n", n)
		fmt.Printf("in test case x4_out=%d\n", x)
	}(x)
	return 9
}

func main() {

	fmt.Println("start")
	x1 := test1()
	fmt.Printf("in main x1=%d\n", *x1)

	x2 := test2()
	fmt.Printf("in main x2=%d\n", x2)

	x3 := test3()
	fmt.Printf("in main x3=%d\n", x3)

	x4 := test4()
	fmt.Printf("in main x4=%d\n", x4)

	fmt.Println("end")
}
