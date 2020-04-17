package main

import (
	"fmt"
)

type People struct{}

func (p *People) ShowA() {
	fmt.Println("ShowA")
	p.ShowB()
}

func (p *People) ShowB() {
	fmt.Println("ShowB")
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	fmt.Println("Show Teacher")
}


func main() {
	t := Teacher{}
	t.ShowA()
}
