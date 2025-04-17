package eventstream

import "context"

type Producer interface {
	Produce(ctx context.Context, event Event) error
	Close(ctx context.Context) error
}
