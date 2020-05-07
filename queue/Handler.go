package queue

import (
	"context"
)

type Handler interface {
	Handle(ctx context.Context) context.Context
	Store(ctx context.Context, m interface{}) context.Context
	Key() string
}

