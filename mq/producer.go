package mq

import "context"

type Producer interface {
	Push(ctx context.Context, topic string, msg interface{}) error
	Close() error
}
