// Package interview 面试相关的题
// 多协程查询切片
// https://github.com/lifei6671/interview-go/blob/master/question/q017.md
// 假设有一个超长的切片，切片的元素类型为int，切片中的元素为乱序排序。限时5秒，使用多个goroutine查找切片中是否存在给定的值，在查找到目标值或者超时后立刻结束所有goroutine的执行。
// 比如，切片 [23,32,78,43,76,65,345,762,......915,86]，查找目标值为 345 ，如果切片中存在，则目标值输出"Found it!"并立即取消仍在执行查询任务的goroutine。
// 如果在超时时间未查到目标值程序，则输出"Timeout！Not Found"，同时立即取消仍在执行的查找任务的goroutine。
package interview

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func i110902(source []int, target int, timeout int, workers int) (err error) {

	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	sign := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())

	l := len(source)
	cutLen := l / workers

	wg := sync.WaitGroup{}
	wg.Add(1 + workers)
	go func() {
		defer wg.Done()
		select {
		case <-timer.C:
			fmt.Println("Timeout! Not Found")
		case <-sign:
			fmt.Println("Found it!")
		}
		cancel()
	}()
	// create work
	for i := 0; i < workers; i++ {
		go func(sub []int) {
			defer wg.Done()
			for _, v := range sub {
				select {
				case <-ctx.Done():
					//结束查找，等待时间到了
					return
				default:
				}
				time.Sleep(1000)
				if v == target {
					sign <- true
					return
				}
			}
		}(source[i*cutLen : (i+1)*cutLen])
	}
	wg.Wait()
	return
}
