package mux

type IMuxRecipient[TElem any] any

// A IMutMuxRecipient is an entity that can subscribe to a multiplexer (mux). Once subscribed, its Process(...) method is
// invoked for each item sent to the mux
type IMutMuxRecipient[TElem any] interface {
	IMuxRecipient[TElem]

	Process(data TElem)
}

// The base structure for a IMutMuxRecipient. Include it into your mux recipients
type MuxRecipient[TElem any] struct {
	Errs chan error
}
