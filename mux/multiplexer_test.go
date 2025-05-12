package mux

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMultiplexerWithCanceledContext(t *testing.T) {
	assert := assert.New(t)

	ctx, ctxCancel := context.WithCancelCause(t.Context())
	ctxCancel(nil)

	ch := make(chan int, 5)

	mux := CreateMultiplexer[int]()

	initErr := mux.Run(ctx, ch)

	assert.NotNil(initErr)
	assert.Contains(initErr.Error(), "abort")
}

type ExampleMuxRecipient struct {
	Base MuxRecipient[int]

	Output chan int
}

func CreateExampleMuxRecipient() ExampleMuxRecipient {
	resps := make(chan int, 1)
	errs := make(chan error, 1)

	return ExampleMuxRecipient{
		Output: resps,
		Base: MuxRecipient[int]{
			Errs: errs,
		},
	}
}

func (p *ExampleMuxRecipient) Process(item int) {
	p.Output <- item * 10
}

func TestMultiplexerMulticastProperties(t *testing.T) {
	assert := assert.New(t)

	ctx, ctxCancel := context.WithCancelCause(t.Context())
	defer ctxCancel(nil)

	ch := make(chan int, 5)

	mux := CreateMultiplexer[int]()

	initErr := mux.Run(ctx, ch)
	assert.Nil(initErr)

	recp := CreateExampleMuxRecipient()

	mux.Subscribe(&recp)

	inputItem := 3
	expectedResult := 30

	ch <- inputItem

	select {
	case actualItem := <-recp.Output:
		assert.Equal(expectedResult, actualItem)

	case <-time.After(1 * time.Millisecond):
		assert.Fail("took too long to multicast using mux")
	}
}
