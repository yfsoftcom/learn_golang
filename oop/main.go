/**


**/
package main

import "fmt"

type BaseType string

func (b BaseType) sayHi() {
	fmt.Println("greet from ", b)
}

type InheritTypeA BaseType

type InheritTypeB struct {
	BaseType
}

func main() {
	itA := InheritTypeA("it A")
	// itA.sayHi() // compile error: the inheritType dont inherit the sayhi method.

	itB := InheritTypeB("it B")
	itB.sayHi() // here is ok

	fmt.Println("hello golang!")
}
