//go:build !test

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SilvioRispoliAtElectroluxDotCom/chanswitch/mux"
)

type (
	happyPrinter struct {
		base mux.MuxRecipient[int]
	}

	sassyPrinter struct {
		base mux.MuxRecipient[int]
	}
)

func (happyPrinter) Process(item int) {
	fmt.Printf("Hey, a beautiful item: %d\n", item)
}

func (sassyPrinter) Process(item int) {
	fmt.Printf("Oh, a item or whatever.. %d\n", item)
}

func main() {
	ctx, ctxCancel := context.WithCancelCause(context.Background())
	defer ctxCancel(nil)

	inlet := make(chan int, 5)

	hp := happyPrinter{}
	sp := sassyPrinter{}

	muxer := mux.CreateMultiplexer[int]()

	muxer.Run(ctx, inlet)
	muxer.Subscribe(hp)
	muxer.Subscribe(sp)

	inlet <- 1

	time.Sleep(1 * time.Second)

	inlet <- 8

	time.Sleep(1 * time.Second)

	return
}
