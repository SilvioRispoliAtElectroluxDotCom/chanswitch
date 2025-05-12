package demux

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeMultiplexerWithCanceledContext(t *testing.T) {
	assert := assert.New(t)

	ctx, ctxCancel := context.WithCancelCause(t.Context())
	ctxCancel(nil)

	demux := CreateDeMultiplexer[int]()

	_, initErr := demux.Run(ctx)

	assert.NotNil(initErr)
	assert.Contains(initErr.Error(), "abort")
}

func TestDeMultiplexerHopper(t *testing.T) {
	assert := assert.New(t)

	ctx, ctxCancel := context.WithCancelCause(t.Context())
	defer ctxCancel(nil)

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	demux := CreateDeMultiplexer[int]()

	collected, initErr := demux.Run(ctx)

	assert.Nil(initErr)

	demux.Subscribe(ch1)
	demux.Subscribe(ch2)

	firstExpected := 5
	secondExpected := 10

	ch1 <- firstExpected
	ch2 <- secondExpected

	firstActual := <-collected
	secondActual := <-collected

	assert.Equal(firstExpected, firstActual)
	assert.Equal(secondExpected, secondActual)
}
