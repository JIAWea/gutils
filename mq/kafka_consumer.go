package mq

import (
	"context"
	"errors"
	"io"

	kg "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaGroup struct {
	kReader *kg.Reader
	logger  *zap.Logger
	Discover
}

func NewKafkaGroup(conf *Conf, logger zap.Logger, discover Discover) *KafkaGroup {
	rc := kg.ReaderConfig{
		Brokers:     conf.Addrs,
		GroupID:     conf.GroupId,
		GroupTopics: conf.Topics,
		StartOffset: kg.LastOffset,
	}
	g := &KafkaGroup{
		kReader:  kg.NewReader(rc),
		logger:   zap.L(),
		Discover: discover,
	}
	return g
}

func (k *KafkaGroup) Close() error {
	return k.kReader.Close()
}

func (k *KafkaGroup) Consume() {
	go func() {
		k.logger.Info("msg", zap.String("id", k.kReader.Config().GroupID))
		for {
			ctx := context.Background()
			msg, err := k.kReader.FetchMessage(ctx)
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrClosedPipe) {
				return
			}
			if err != nil {
				k.logger.Error("err", zap.Error(err))
				continue
			}
			handler, err := k.Discover.GetHandler(msg.Topic)
			if err != nil {
				k.logger.Error("err", zap.Error(err))
				continue
			}

			err = handler(ctx, msg.Topic, msg.Key, msg.Value)
			if err != nil {
				continue
			}
			err = k.kReader.CommitMessages(ctx, msg)
			if err != nil {
				k.logger.Error("err", zap.Error(err))
			}
		}
	}()

}
