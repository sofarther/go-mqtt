package main

import (
	"context"
	"mqtt2/broadcast"
	"mqtt2/queue"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {



	ctx ,cancel := context.WithCancel(context.Background())

	//defer cancel()

	ch := make(chan os.Signal)
	signal.Notify(ch,os.Interrupt,syscall.SIGTERM)

	//queue.Register(func(m interface{}) {
	//	//fmt.Printf("receive message: %s\n", m)
	//	broadcast.Input(m.(string))
	//})

	// 保证 调用顺序：
	// 先 注册 活动的 topic ，并为 每个 topic 开启 协程监听
	// 然后 注册 消息处理回调函数， 并 开启协程 等待 消息
	// 最后 开启协程 获取消息
	// 否则 应用启动时， 先获取的消息 将无法 被广播出去
	broadcast.Start(ctx)

	queue.Register(queue.NewLoggerHandler())

	queue.Register(broadcast.NewDecodeHandler())
	queue.Register(broadcast.NewPublishHandler())
	queue.HandleMessage(ctx)


	queue.Start(ctx)


	 <- ch

	cancel()

	time.Sleep(5 * time.Second)
}
