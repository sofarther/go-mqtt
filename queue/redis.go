package queue

import (
	"context"
	"fmt"
	"log"
	"mqtt2/initialize"
	"time"
)

const queueName  = "go_queue"


func getMessage() string {
	cmd := initialize.Redis().BRPop(10 * time.Second, queueName)

	result, err := cmd.Result()
	if err != nil {
		fmt.Println("get message error: " , err.Error())
		return ""
	}

	// fmt.Println(result)
	return result[1]
}

func start(ctx context.Context)  {
	go func() {

		defer log.Println("start done")
		for {
			select {
			case <- ctx.Done():
				return
			default:
				m := getMessage()
				fmt.Println("m:", m)
				if m == "" {
					time.Sleep(1 * time.Second)
					continue
				}
				messageChan <- m
			}

		}
	}()

}
