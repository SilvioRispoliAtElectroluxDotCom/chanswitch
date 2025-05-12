package demux

import (
	"context"
	"fmt"
)

// A DeMultiplexer collects items from many channels and funnels them into one
type DeMultiplexer[T any] struct {
	sources []chan T
}

func CreateDeMultiplexer[T any]() DeMultiplexer[T] {
	return DeMultiplexer[T]{
		sources: []chan T{},
	}
}

// Subscribes a channel to the demultiplexer. Items will be pulled asynchronously during the Run(...) goroutine
func (ec *DeMultiplexer[T]) Subscribe(source chan T) {
	ec.sources = append(ec.sources, source)
}

// Spins a goroutine that pulls items from the subscribed channels and puts them onto the returned channel
func (ec *DeMultiplexer[T]) Run(ctx context.Context) (chan T, error) {
	collected := make(chan T, 1)

	if err := ctx.Err(); err != nil {
		close(collected)
		return collected, fmt.Errorf("abort DeMultiplexer due to error: %w", context.Cause(ctx))
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				for _, ch := range ec.sources {
					select {
					case el := <-ch:
						collected <- el
					default:

					}
				}
			}
		}
	}()

	return collected, nil
}
