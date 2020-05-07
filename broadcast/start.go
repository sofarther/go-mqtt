package broadcast

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func registerAliveTopic(ctx context.Context, client mqtt.Client)  {

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

			go run(client, topicCh[v])
		}
	}

	for k, v := range topicCh {
		if _, ok := aliveTopicMap[k]; !ok {
			v.cancel()
			delete(topicCh,k)
		}
	}

}

func run(client mqtt.Client, topic *topicChannel)  {

	for  {
		select {
		case m := <-topic.ch:
			if token := client.Publish(topic.topic,1,false,m); token.Wait() && token.Error() != nil {
				log.Fatalf("publish error: %s \n", token.Error())
			}
		case <-topic.ctx.Done():
			log.Println(topic.topic , " done")
			return

		}
	}
}


func Start(ctx context.Context)  {

	client := newClient()
	// 启动时 注册 活跃的 topic
	registerAliveTopic(ctx, client)
   // 定时 处理 活跃的 topic
	go func() {
		for  {
			select {
			case <- time.After(5 * time.Second):
				registerAliveTopic(ctx, client)
			case <- ctx.Done():
				return


			}
		}
	}()


	// go registerAliveTopic(ctx, client)
}

