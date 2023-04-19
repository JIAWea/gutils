package mq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"go.uber.org/zap"
)

var (
	ErrHandlerNotFound = errors.New("topic handler not found")
)

type Conf struct {
	Driver  string
	Addrs   []string
	Topics  []string
	GroupId string
}

type Consumer interface {
	Consume()
	Close() error
}

type Discover interface {
	GetHandler(topic string) (ConsumeHandle, error)
}

type Register interface {
	RegisterHandler(topic string, h ConsumeHandle)
}

type ConsumeHandle func(ctx context.Context, topic string, k, v []byte) error

type Handlers struct {
	handlers map[string]ConsumeHandle
	logger   *zap.Logger
}

var (
	hm   *Handlers
	hone sync.Once
)

func NewHandlerManager(logger log.Logger) *Handlers {
	hone.Do(func() {
		hm = &Handlers{
			handlers: map[string]ConsumeHandle{},
			logger:   zap.L(),
		}
	})
	return hm
}

func (a *Handlers) RegisterHandler(topic string, h ConsumeHandle) {
	a.handlers[topic] = h
}

func (a *Handlers) GetHandler(topic string) (ConsumeHandle, error) {
	h, ok := a.handlers[topic]
	if !ok {
		a.logger.Error(fmt.Sprintf("%v", topic))
		return nil, ErrHandlerNotFound
	}
	return h, nil
}

func (a *Handlers) Topics() []string {
	topics := make([]string, 0, len(a.handlers))
	for t := range a.handlers {
		topics = append(topics, t)
	}
	return topics
}
