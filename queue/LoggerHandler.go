package queue

import (
	"context"
	"log"
)

type LoggerHandler struct {
	key string
}

func (h *LoggerHandler) Handle(ctx context.Context)  context.Context {
	m := ctx.Value(GetMessageOriginKey())
	log.Printf("log message:%s\n",m)

	return ctx
}
func (h *LoggerHandler)Store(ctx context.Context, m interface{})  context.Context{

	return ctx
}
func (h *LoggerHandler) Key()  string{
	return h.key
}


func NewLoggerHandler() *LoggerHandler {
	return &LoggerHandler{key:"logger"}
}

