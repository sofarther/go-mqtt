package broadcast

import (
	"context"
	"mqtt2/queue"
)

type decoderHandler struct {
	key string
}

type publishHandler struct {
	key string

}

var decoderH *decoderHandler
var publishH *publishHandler

func (p *publishHandler) Handle(ctx context.Context) context.Context {
	m := ctx.Value(decoderH.Key()).(Message)


	ch, ok := topicCh[m.Topic]
	if !ok {
		return  ctx
	}

	ch.ch <- m.Payload
	// m.Topic

	return  ctx
}

func (p *publishHandler) Store(ctx context.Context, m interface{}) context.Context {
	return ctx
}

func (p *publishHandler) Key() string {
	return p.key
}

func (d *decoderHandler) Handle(ctx context.Context) context.Context {
	 m := queue.GetOriginMessage(ctx).(string)

	 message ,err := Decode(m)
	 if err != nil {
	 	panic(err)
	 }

	 return d.Store(ctx, message)
}

func (d *decoderHandler) Store(ctx context.Context, m interface{}) context.Context {
	return context.WithValue(ctx, d.key, m)
}

func (d *decoderHandler) Key() string {
	return d.key
}

func NewDecodeHandler() queue.Handler {
	if decoderH == nil {
		decoderH = &decoderHandler{key:"decoder"}
	}

	return decoderH
}

func NewPublishHandler() queue.Handler {
	if publishH == nil {
		publishH = &publishHandler{key:"push"}
	}

	return publishH
}

