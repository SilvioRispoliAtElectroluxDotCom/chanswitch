package mux

import (
	"context"
	"fmt"
)

// A Multiplexer takes items from a channel and multicast them asynchronously to all the subscribers
type Multiplexer[T any] struct {
	recipients []IMutMuxRecipient[T]
}

func CreateMultiplexer[T any]() Multiplexer[T] {
	return Multiplexer[T]{}
}

// Subscribe a recipient to this mux
func (mux *Multiplexer[T]) Subscribe(recipient IMutMuxRecipient[T]) {
	mux.recipients = append(mux.recipients, recipient)
}

// Spins a goroutine that multicasts incoming items from the channel to all subscribed recipients
func (mux *Multiplexer[T]) Run(ctx context.Context, ch chan T) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("abort Multiplexer due to error: %w", context.Cause(ctx))
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case elem := <-ch:
				for _, recipient := range mux.recipients {
					recipient.Process(elem)
				}
			}
		}
	}()

	return nil
}
