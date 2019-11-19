package main

import "fmt"

/**
defer 相关内容

1. 运行时机
	1.1 函数内的 defer 是否会在函数外部也执行？ [并不会]
	1.2 匿名函数中的defer何时执行？ 匿名函数中的 defer 会在该函数生命周期里完成，不会被外层的函数干扰
2. 执行顺序
	2.1 先定义的 defer 后执行，遵循堆栈的运行模型
	2.2 是否无论在函数的什么位置都必定会执行？  在 panic 后的 defer 是无法正常执行的
3. 与 panic 配合

*/

func foo() {
	defer fmt.Println("foo2")
	for i := 0; i < 3; i++ {
		fmt.Println(i)
		(func() {
			defer fmt.Println("sub", i)
		})()
	}

	defer fmt.Println("foo1")
}

func main() {
	fmt.Println("start")
	foo()
	fmt.Println("end")
}
