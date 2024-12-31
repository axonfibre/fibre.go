package eventticker

import (
	"github.com/axonfibre/fibre.go/core/index"
	"github.com/axonfibre/fibre.go/runtime/event"
)

// Events represents events happening on a EventTicker.
type Events[I index.Type, T index.IndexedID[I]] struct {
	// Request is an event that is triggered when the requester wants to request the given entity.
	Tick *event.Event1[T]

	// RequestQueued is an event that is triggered when a new request is started.
	TickerStarted *event.Event1[T]

	// RequestStopped is an event that is triggered when a request is stopped.
	TickerStopped *event.Event1[T]

	// RequestFailed is an event that is triggered when a request is stopped after too many attempts.
	TickerFailed *event.Event1[T]

	event.Group[Events[I, T], *Events[I, T]]
}

// NewEvents contains the constructor of the Events object (it is generated by a generic factory).
func NewEvents[I index.Type, T index.IndexedID[I]](linkedEvents ...*Events[I, T]) (newEvents *Events[I, T]) {
	return event.CreateGroupConstructor(func() *Events[I, T] {
		return &Events[I, T]{
			Tick:          event.New1[T](),
			TickerStarted: event.New1[T](),
			TickerStopped: event.New1[T](),
			TickerFailed:  event.New1[T](),
		}
	})(linkedEvents...)
}
