/**
获取目录下的所有文件的大小总和
方案1,使用递归的方式获取
方案2,使用goroutine并发获取
方案3,使用channel控制并发上限10

output:

dir: %s
total: %d
**/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var totalSize int64

func main() {

	dir := "/home/wangfan/Desktop"

	var total int64
	var err error
	total, err = grab2(dir)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("dir: %s\ntotal: %d", dir, total)
}

func grab2(dir string) (int64, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return -1, err
	}
	var t int64
	for _, f := range files {
		if f.IsDir() {
			// it's a sub dir
			total, anotherErr := grab1(dir + string(os.PathSeparator) + f.Name())
			if anotherErr != nil {
				return -1, anotherErr
			}
			t += total
		} else {
			t += f.Size()
		}
	}
	return t, nil
}
