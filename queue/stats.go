package queue

import (
	"context"
	"encoding"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	keyStatsIncr = "queue:stats:incr"
)

type StatsIncrEvent struct {
	LinkID string            `redis:"linkId"`
	Labels map[string]string `redis:"labels"`
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *StatsIncrEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s *StatsIncrEvent) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}

var (
	_ encoding.BinaryMarshaler   = (*StatsIncrEvent)(nil)
	_ encoding.BinaryUnmarshaler = (*StatsIncrEvent)(nil)
)

type StatsQueue struct {
	redis *redis.Client

	key string
}

func (q *StatsQueue) Enqueue(ctx context.Context, event StatsIncrEvent) error {
	return q.redis.RPush(ctx, q.key, &event).Err()
}

func (q *StatsQueue) Dequeue(ctx context.Context) (StatsIncrEvent, error) {
	event := StatsIncrEvent{}

	vals, err := q.redis.BLPop(ctx, time.Second, q.key).Result()
	if err == redis.Nil {
		return StatsIncrEvent{}, ErrEmptyQueue
	}

	if err != nil {
		return StatsIncrEvent{}, err
	}

	if len(vals) != 2 {
		return StatsIncrEvent{}, ErrEmptyQueue
	}

	if err := event.UnmarshalBinary([]byte(vals[1])); err != nil {
		return StatsIncrEvent{}, err
	}

	return event, nil
}

func NewStatsQueue(redis *redis.Client) *StatsQueue {
	if redis == nil {
		panic("redis client is nil")
	}

	return &StatsQueue{
		redis: redis,

		key: keyStatsIncr,
	}
}
