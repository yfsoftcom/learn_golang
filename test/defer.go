//
/*
Guess Output:
20 0 1 1
2 0 1 1
10 0 1 1
1 0 1 1

True Output:
10 1 2 3
20 0 2 2
2 0 2 2
1 1 3 4

defer 函数虽然在方法最后执行，但是其参数列表中的方法会按照顺序执行
*/
//

package main

import "fmt"

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func main() {
	a := 1
	b := 2
	defer calc("1", a, calc("10", a, b))
	a = 0
	defer calc("2", a, calc("20", a, b))
	b = 1
}
