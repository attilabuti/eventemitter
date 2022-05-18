package eventemitter

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	ErrEmptyName      = errors.New("Event name cannot be empty")
	ErrNotAFunction   = errors.New("Callback must be a function")
	ErrEventExists    = errors.New("Event already exists")
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
func (e *Emitter) AddListener(eventName string, listener any) error {
	if len(eventName) == 0 {
		return ErrEmptyName
	}

	if listener == nil || reflect.TypeOf(listener).Kind() != reflect.Func {
		return ErrNotAFunction
	}

	if listeners, ok := e.listeners.Load(eventName); ok {
		ptr := reflect.ValueOf(listener).Pointer()

		for _, handler := range listeners.([]any) {
			if reflect.ValueOf(handler).Pointer() == ptr {
				return ErrEventExists
			}
		}

		listeners = append(listeners.([]any), listener)
		e.listeners.Store(eventName, listeners)
	} else {
		e.listeners.Store(eventName, []any{listener})
	}

	return nil
}

// On is an alias for .AddListener(eventName, listener).
func (e *Emitter) On(eventName string, listener any) error {
	return e.AddListener(eventName, listener)
}

// RemoveListener removes the listener for the specified event.
// Returns an error if the event does not exist, or the listener is not a function.
func (e *Emitter) RemoveListener(eventName string, listener any) error {
	if len(eventName) == 0 {
		return ErrEmptyName
	}

	if listener == nil || reflect.TypeOf(listener).Kind() != reflect.Func {
		return ErrNotAFunction
	}

	if listeners, ok := e.listeners.Load(eventName); ok {
		ptr := reflect.ValueOf(listener).Pointer()

		for i, handler := range listeners.([]any) {
			if reflect.ValueOf(handler).Pointer() == ptr {
				if len(listeners.([]any)) == 1 {
					e.listeners.Delete(eventName)
				} else {
					listeners = append(listeners.([]any)[:i], listeners.([]any)[i+1:]...)
					e.listeners.Store(eventName, listeners)
				}

				return nil
			}
		}
	}

	return ErrEventNotExists
}

// Off is an alias for .RemoveListener(eventName, listener).
func (e *Emitter) Off(eventName string, listener any) error {
	return e.RemoveListener(eventName, listener)
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
func (e *Emitter) Emit(eventName string, arguments ...any) error {
	return e.emit(eventName, arguments, false)
}

// EmitSync emits an event synchronously with the specified arguments.
// Returns an error if the event does not exist.
func (e *Emitter) EmitSync(eventName string, arguments ...any) error {
	return e.emit(eventName, arguments, true)
}

func (e *Emitter) emit(eventName string, arguments []any, sync bool) error {
	if len(eventName) == 0 {
		return ErrEmptyName
	}

	if listeners, ok := e.listeners.Load(eventName); ok {
		args := make([]reflect.Value, 0)
		for _, param := range arguments {
			args = append(args, reflect.ValueOf(param))
		}

		for _, listener := range listeners.([]any) {
			handler := reflect.ValueOf(listener)
			if err := e.checkArguments(eventName, handler.Type(), args); err != nil {
				panic(err)
			}

			if sync {
				handler.Call(args)
			} else {
				go handler.Call(args)
			}
		}

		return nil
	}

	return ErrEventNotExists
}

func (e *Emitter) checkArguments(eventName string, fnType reflect.Type, args []reflect.Value) error {
	numIn := fnType.NumIn()

	// Check arguments length.
	if fnType.IsVariadic() {
		if (numIn - 1) > len(args) {
			return fmt.Errorf("Not enough arguments. Event %s expected at least %d arguments, got %d.", eventName, numIn-1, len(args))
		}
	} else if numIn != len(args) {
		return fmt.Errorf("Wrong number of arguments. Event %s expected %d arguments, got %d.", eventName, numIn, len(args))
	}

	// Check arguments type.
	for i := 0; i < numIn; i++ {
		if fnType.IsVariadic() && i == (numIn-1) {
			variadicArgs := args[i:]

			if len(variadicArgs) > 0 {
				variadicType := fnType.In(i).Elem()

				for j := 0; j < len(variadicArgs); j++ {
					if !variadicArgs[j].Type().AssignableTo(variadicType) {
						return fmt.Errorf("Wrong argument type. Event %s expected argument %d to be %s, got %s.", eventName, i+1, variadicType, variadicArgs[j].Type())
					}
				}
			}
		} else if !args[i].Type().AssignableTo(fnType.In(i)) {
			return fmt.Errorf("Wrong argument type. Event %s expected argument %d to be %s, got %s.", eventName, i+1, fnType.In(i), args[i].Type())
		}
	}

	return nil
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

// Listeners returns a slice of functions registered to the specified event.
// Returns an error if the event does not exist.
func (e *Emitter) Listeners(eventName string) ([]any, error) {
	return e.getListeners(eventName)
}

// ListenersCount returns the number of listeners for the specified event.
// Returns an error if the event does not exist.
func (e *Emitter) ListenerCount(eventName string) (int, error) {
	listeners, err := e.getListeners(eventName)

	if err != nil {
		return 0, err
	}

	return len(listeners), nil
}

func (e *Emitter) getListeners(eventName string) ([]any, error) {
	if len(eventName) == 0 {
		return nil, ErrEmptyName
	}

	if listeners, ok := e.listeners.Load(eventName); ok {
		return listeners.([]any), nil
	}

	return nil, ErrEventNotExists
}
