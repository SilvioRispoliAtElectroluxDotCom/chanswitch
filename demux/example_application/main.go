//go:build !test

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SilvioRispoliAtElectroluxDotCom/chanswitch/demux"
)

func main() {
	ctx, ctxCancel := context.WithCancelCause(context.Background())
	defer ctxCancel(nil)

	demux := demux.CreateDeMultiplexer[int]()
	ints, err := demux.Run(ctx)

	if err != nil {
		panic("Something went wrong while initializing the demux")
	}

	// These in a real application would usually be channels returned by some goroutine
	ch1 := make(chan int, 2)
	ch2 := make(chan int, 5)

	demux.Subscribe(ch1)
	demux.Subscribe(ch2)

	// These would happen withing the bounds of the goroutine at their own leisure
	ch1 <- 1
	ch2 <- 3
	ch1 <- 5
	ch2 <- 7
	ch2 <- 9

	// Simulating expensive operation here
	time.Sleep(2 * time.Second)

	for {
		select {
		case i := <-ints:
			fmt.Printf("Got item: %d\n", i)

		default:
			return
		}
	}
}
