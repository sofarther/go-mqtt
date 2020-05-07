package queue

import (
	"context"
	"errors"
	"sync"
)

const messageOriginKey  = "origin"
var messageChan chan string
func init()  {
	messageChan = make(chan string)
}


// type handler func(ctx context.Context)
//type handler struct {
//	name string
//}

type handlerSlice []Handler
type handlerKeySlice []string


var handlers handlerSlice
var handlerKey handlerKeySlice

var mutex sync.Mutex

func Register( h Handler) (err error) {
	mutex.Lock()
	defer mutex.Unlock()

	key := h.Key()
	if key == "" {
		return errors.New("key cannot empty")
	}
	if key == messageOriginKey {
		return errors.New("key has exists")
	}
	for _, v := range handlerKey {
		if v == key {
			return errors.New("key has exists")
		}
	}

	handlerKey = append(handlerKey, key)
	handlers = append(handlers, h)

	return nil
}

func HandleMessage(ctx context.Context)  {
	go func() {
		//ctx, _ = context.WithCancel(ctx)

		for  {
			select {
			case <- ctx.Done():
				return
			case m := <- messageChan:
				valCtx := context.WithValue(ctx, messageOriginKey, m )
				for _, h := range handlers {
					valCtx = h.Handle(valCtx)

				}
			}
		}
	}()

}

func Start(ctx context.Context)  {
	//ctx, _ = context.WithCancel(ctx)

	start(ctx)
}

func GetMessageOriginKey()  string{
	return messageOriginKey
}

func GetOriginMessage(ctx context.Context) interface{} {
	return ctx.Value(messageOriginKey)
}


