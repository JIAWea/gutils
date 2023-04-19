package mq

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	kg "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	kwriter *kg.Writer
	logger  *zap.Logger
}

func NewKafkaProducer(conf *Conf, logger *zap.Logger) *KafkaProducer {
	kw := &kg.Writer{
		Addr:                   kg.TCP(conf.Addrs...),
		Async:                  true,
		Completion:             nil,
		Compression:            0,
		Logger:                 nil,
		AllowAutoTopicCreation: true,
	}
	t := &KafkaProducer{
		kwriter: kw,
		logger:  zap.L(),
	}
	return t
}

func (p *KafkaProducer) traceHeader(ctx context.Context) (kg.Header, error) {
	spanContext := trace.SpanFromContext(ctx).SpanContext()
	bytes, err := spanContext.MarshalJSON()
	if err != nil {
		return kg.Header{}, err
	}
	return kg.Header{
		Key:   spanKey,
		Value: bytes,
	}, nil
}

func (p *KafkaProducer) Push(ctx context.Context, topic string, msg interface{}) error {
	header, err := p.traceHeader(ctx)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	key := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	km := kg.Message{
		Key:     key,
		Topic:   topic,
		Value:   bytes,
		Headers: []kg.Header{header},
	}
	err = p.kwriter.WriteMessages(ctx, km)
	if err != nil {
		p.logger.Error("msg", zap.Error(err))
		return err
	}
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.kwriter.Close()
}
