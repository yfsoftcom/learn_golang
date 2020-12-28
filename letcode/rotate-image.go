// main
// the solution of https://leetcode-cn.com/problems/rotate-image/submissions/
package main

import "fmt"

func rotate(matrix [][]int) {
	SIZE := len(matrix) - 1
	XRange := len(matrix) / 2
	if len(matrix)%2 > 0 {
		XRange++
	}
	for y := 0; y < XRange-1; y++ {
		for x := 0; x < XRange; x++ {
			// 进行旋转
			matrix[x][SIZE-y], matrix[SIZE-x][SIZE-y], matrix[SIZE-x][y], matrix[x][y] = matrix[x][y], matrix[x][SIZE-y], matrix[SIZE-x][SIZE-y], matrix[SIZE-x][y]
		}
	}
}

func main() {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	rotate(matrix)
	fmt.Println("result: \n", matrix)
}
