package queue

import (
	"context"
	"time"
)

type Subscriber interface {
	Init(opts ...Option) error
	Handle(ctx context.Context, numWorkers int, reconnectTimeout time.Duration)
}
