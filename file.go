package main

import (
	"fmt"
	"time"
)

/*
*
关闭一个 channel 可以使用内置函数 close()
*/
func main() {
	// 创建一个 channel
	ch := make(chan int)

	// 启动一个 goroutine 向 channel 发送数据
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(time.Millisecond * 100)
		}
		// 关闭 channel
		close(ch)
	}()

	// 从 channel 中接收数据
	for {
		// 使用 ok 来判断 channel 是否已关闭
		value, ok := <-ch
		if !ok {
			fmt.Println("Channel closed.")
			break
		}
		fmt.Println("Received:", value)
	}
}
