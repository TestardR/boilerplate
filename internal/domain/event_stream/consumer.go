package eventstream

import "context"

type EventHandler func(ctx context.Context, event Event) error

type Consumer interface {
	Consume(ctx context.Context, handler EventHandler) error
	Close(ctx context.Context) error
}
