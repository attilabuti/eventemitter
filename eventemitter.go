package eventemitter

import (
	"errors"
	"reflect"
	"sync"
)

var (
	// ErrEmptyName is returned when the event name is empty.
	ErrEmptyName = errors.New("Event name cannot be empty")

	// ErrNotAFunction is returned when the callback is not a function.
	ErrNotAFunction = errors.New("Callback must be a function")

	// ErrEventExists is returned when the event already exists.
	ErrEventExists = errors.New("Event already exists")

	// ErrEventNotExists is returned when the event does not exist.
	ErrEventNotExists = errors.New("Event does not exist")
)

type Emitter struct {
	listeners sync.Map
}

// New returns a new event emitter.
func New() *Emitter {
	return &Emitter{}
}

// AddListener adds a listener for the specified event.
// Returns an error if the event already exists, or the listener is not a function.
func (e *Emitter) AddListener(eventName string, listener any) (err error) {
	if len(eventName) == 0 {
		return ErrEmptyName
	}

	if listener == nil || reflect.TypeOf(listener).Kind() != reflect.Func {
		return ErrNotAFunction
	}

	_, ok := e.listeners.LoadOrStore(eventName, listener)
	if ok {
		return ErrEventExists
	}

	return nil
}

// On is an alias for .AddListener(eventName, listener).
func (e *Emitter) On(eventName string, listener any) (err error) {
	return e.AddListener(eventName, listener)
}

// RemoveListener removes the listener for the specified event.
// Returns an error if the event does not exist.
func (e *Emitter) RemoveListener(eventName string) (err error) {
	if len(eventName) == 0 {
		return ErrEmptyName
	}

	_, ok := e.listeners.LoadAndDelete(eventName)
	if !ok {
		return ErrEventNotExists
	}

	return nil
}

// Off is an alias for .RemoveListener(eventName).
func (e *Emitter) Off(eventName string) (err error) {
	return e.RemoveListener(eventName)
}

// RemoveAllListeners removes all listeners, or those of the specified eventName.
func (e *Emitter) RemoveAllListeners(eventName ...string) {
	if len(eventName) == 0 {
		eventName = e.EventNames()
	}

	for _, event := range eventName {
		e.listeners.Delete(event)
	}
}

// Clear is an alias for .RemoveAllListeners(eventName).
func (e *Emitter) Clear(eventName ...string) {
	e.RemoveAllListeners(eventName...)
}

// Emit emits an event asynchronously with the specified arguments.
// Returns an error if the event does not exist.
func (e *Emitter) Emit(eventName string, arguments ...any) (err error) {
	return e.emit(eventName, arguments, false)
}

// EmitSync emits an event synchronously with the specified arguments.
// Returns an error if the event does not exist.
func (e *Emitter) EmitSync(eventName string, arguments ...any) (err error) {
	return e.emit(eventName, arguments, true)
}

func (e *Emitter) emit(eventName string, arguments []any, sync bool) error {
	if len(eventName) == 0 {
		return ErrEmptyName
	}

	if listener, ok := e.listeners.Load(eventName); ok {
		args := make([]reflect.Value, 0)
		for _, arg := range arguments {
			args = append(args, reflect.ValueOf(arg))
		}

		if sync {
			reflect.ValueOf(listener).Call(args)
		} else {
			go reflect.ValueOf(listener).Call(args)
		}

		return nil
	}

	return ErrEventNotExists
}

// EventNames returns a slice of strings listing the events for which the emitter
// has registered listeners.
func (e *Emitter) EventNames() []string {
	var names []string

	e.listeners.Range(func(eventName, listener any) bool {
		names = append(names, eventName.(string))
		return true
	})

	return names
}
