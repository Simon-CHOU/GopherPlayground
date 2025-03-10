//Go可以使用 context包中的 Context 类型来控制 Goroutine 的生命周期和取消

package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 创建一个带有取消功能的 context
	ctx, cancel := context.WithCancel(context.Background())

	// 启动一个 goroutine 执行耗时操作
	go worker(ctx)

	// 模拟主程序运行 3 秒
	time.Sleep(3 * time.Second)

	// 取消 context
	cancel()

	// 等待一段时间，确保 worker goroutine 有足够时间退出
	time.Sleep(1 * time.Second)
	fmt.Println("主程序退出")
}

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// context 被取消时，退出 goroutine
			fmt.Println("Worker: 收到取消信号，正在退出...")
			return
		default:
			// 模拟工作
			fmt.Println("Worker: 正在工作...")
			time.Sleep(1 * time.Second)
		}
	}
}
