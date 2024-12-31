package reactive

// event is the default implementation of the Event interface.
type event struct {
	Variable[bool]
}

// newEvent creates a new event.
func newEvent() *event {
	return &event{
		Variable: NewVariable[bool](func(currentValue bool, newValue bool) bool {
			// make sure that the value will always be true once it was set to true
			return currentValue || newValue
		}),
	}
}

// Trigger triggers the event and returns true if the event was triggered for the first time.
func (e *event) Trigger() bool {
	return !e.Set(true)
}

// WasTriggered returns true if the event was triggered.
func (e *event) WasTriggered() bool {
	return e.Get()
}

// OnTrigger registers a callback that is executed when the event is triggered.
func (e *event) OnTrigger(handler func()) (unsubscribe func()) {
	return e.OnUpdate(func(_, _ bool) {
		handler()
	})
}
