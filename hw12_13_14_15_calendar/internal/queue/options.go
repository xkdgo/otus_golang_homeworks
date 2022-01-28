package queue

import (
	"context"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
)

func WithLogger(logg *logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logg
	}
}

func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

type Option func(*Options)

type Options struct {
	Logger *logger.Logger
	// Alternative options
	Context context.Context
}
