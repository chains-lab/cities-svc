package publisher

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

type Service struct {
	addr   string
	writer *kafka.Writer
}

func New(addr string) *Service {
	return &Service{
		addr: addr,
		writer: &kafka.Writer{
			Addr:         kafka.TCP(addr),
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
			Compression:  kafka.Snappy,
			BatchTimeout: 50 * time.Millisecond,
		},
	}
}

func (s *Service) Close() error {
	if s.writer != nil {
		return s.writer.Close()
	}
	return nil
}

type Envelope interface {
	MarshalJSON() ([]byte, error)
	EventType() string
	EventVersion() string
	EventTime() time.Time
}

func (s *Service) publish(
	ctx context.Context,
	topic, key string,
	envelope Envelope,
	headers ...kafka.Header,
) error {
	body, err := envelope.MarshalJSON()
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: body,
		Time:  envelope.EventTime(),
		Headers: append(headers,
			kafka.Header{Key: "event_type", Value: []byte(envelope.EventType())},
			kafka.Header{Key: "event_version", Value: []byte(envelope.EventVersion())},
			kafka.Header{Key: "content-type", Value: []byte("application/json")},
		),
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	var lastErr error
	backoff := 100 * time.Millisecond
	for i := 0; i < 3; i++ {
		if err = s.writer.WriteMessages(ctx, msg); err == nil {
			return nil
		}
		lastErr = err
		if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			break
		}
		time.Sleep(backoff)
		backoff *= 2
	}

	return lastErr
}
