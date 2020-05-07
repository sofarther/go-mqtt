package broadcast

import (
	"context"
	"log"
	"mqtt2/initialize"
	"sync"
	"time"
)

const aliveTopicSetKey  = "alive:topic"

type topicChannel struct {
	topic string
	ch chan string
	cancel context.CancelFunc
	ctx context.Context
}

type topicChannels map[string]*topicChannel

var topicCh topicChannels

// var messageCh chan string

var mutex sync.Mutex

func init()  {
	// messageCh  = make(chan string)
	topicCh = make(topicChannels)
}

func getAliveTopic() []string {
	redis := initialize.Redis()

	result,err := redis.SMembers(aliveTopicSetKey).Result()

	if err != nil {
		panic(nil)
	}

	return result
}

func registerAliveTopic(ctx context.Context)  {

	aliveTopic := getAliveTopic()

	mutex.Lock()
	defer mutex.Unlock()

	var aliveTopicMap = map[string]bool{}

	for _, v := range aliveTopic {
		aliveTopicMap[v] = true
		if _, ok := topicCh[v]; !ok {
			cancelCtx ,cancel := context.WithCancel(ctx)
			topicCh[v] = &topicChannel{
				topic: v,
				ch:  make(chan string),
				cancel: cancel,
				ctx: cancelCtx,
			}

			go run(topicCh[v])
		}
	}

	for k, v := range topicCh {
		if _, ok := aliveTopicMap[k]; !ok {
			v.cancel()
			delete(topicCh,k)
		}
	}

}

func run(topic *topicChannel)  {

	for  {
		select {
		case m := <-topic.ch:

			publish(topic.topic, m, 1)

		case <-topic.ctx.Done():
			log.Println(topic.topic , " done")
			return

		}
	}
}


func Start(ctx context.Context)  {

	// 启动时 注册 活跃的 topic
	registerAliveTopic(ctx)
   // 定时 处理 活跃的 topic
	go func() {
		for  {
			select {
			case <- time.After(5 * time.Second):
				registerAliveTopic(ctx)
			case <- ctx.Done():
				return


			}
		}
	}()


}

