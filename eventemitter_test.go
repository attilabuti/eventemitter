package eventemitter

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testType struct {
	name string
}

func TestNewEmitter(t *testing.T) {
	emitter := New()
	em := &Emitter{}

	assert.IsType(t, &Emitter{}, emitter)
	assert.Exactly(t, em, emitter)
}

func TestAddListenerOn(t *testing.T) {
	emitter := New()

	// Register an event.
	err := emitter.On("event", func() {})
	if assert.NoError(t, err) {
		_, ok := emitter.listeners.Load("event")
		assert.True(t, ok)
	}

	// Register an event.
	err = emitter.AddListener("add_listener", func() {})
	if assert.NoError(t, err) {
		_, ok := emitter.listeners.Load("add_listener")
		assert.True(t, ok)
	}

	// Register same event twice.
	event := func() {}
	emitter.On("same_event", event)
	err = emitter.On("same_event", event)
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventExists, err)

		_, ok := emitter.listeners.Load("same_event")
		assert.True(t, ok)
	}

	// Empty event name.
	err = emitter.AddListener("", func() {})
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
	}

	// Empty event name.
	err = emitter.On("", func() {})
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
	}

	// Empty event handler.
	err = emitter.AddListener("nil", nil)
	if assert.Error(t, err) {
		assert.Equal(t, ErrNotAFunction, err)
	}

	// Register an event without handler.
	err = emitter.On("nil", "not a function")
	if assert.Error(t, err) {
		assert.Equal(t, ErrNotAFunction, err)
	}
}

func TestRemoveListenerOff(t *testing.T) {
	emitter := New()

	event_1 := func() {}
	event_2 := func() {}
	event_3 := func() {}
	event_3_1 := func() {}

	// Register events.
	emitter.On("event_first", event_1)
	emitter.AddListener("event_second", event_2)
	emitter.On("event_third", event_3)
	emitter.On("event_third", event_3_1)

	// Remove an event.
	err := emitter.Off("event_first", event_1)
	if assert.NoError(t, err) {
		_, ok := emitter.listeners.Load("event_first")
		assert.False(t, ok)
	}

	// Event does not exist.
	err = emitter.RemoveListener("event_third", func() {})
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventNotExists, err)
	}

	// Remove an event.
	listeners, ok := emitter.listeners.Load("event_third")
	assert.True(t, ok)
	assert.Equal(t, 2, len(listeners.([]any)))

	err = emitter.Off("event_third", event_3)
	if assert.NoError(t, err) {
		listeners, ok := emitter.listeners.Load("event_third")
		assert.True(t, ok)
		assert.Equal(t, 1, len(listeners.([]any)))

		err = emitter.Off("event_third", event_3)
		if assert.Error(t, err) {
			assert.Equal(t, ErrEventNotExists, err)
		}
	}

	err = emitter.RemoveListener("event_third", event_3_1)
	if assert.NoError(t, err) {
		_, ok := emitter.listeners.Load("event_third")
		assert.False(t, ok)

		err = emitter.Off("event_third", event_3_1)
		if assert.Error(t, err) {
			assert.Equal(t, ErrEventNotExists, err)
		}
	}

	// Event does not exist.
	err = emitter.RemoveListener("event_fourth", func() {})
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventNotExists, err)
	}

	// Empty event name.
	err = emitter.Off("", func() {})
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
	}

	// Remove an event without handler.
	err = emitter.RemoveListener("event_second", "not a function")
	if assert.Error(t, err) {
		assert.Equal(t, ErrNotAFunction, err)
	}
}

func TestRemoveAllListenersClear(t *testing.T) {
	emitter := New()

	// Register events.
	emitter.On("event_first", func() {})
	emitter.AddListener("event_second", func() {})
	emitter.On("event_third", func() {})

	// Remove multiple events.
	emitter.RemoveAllListeners("event_first", "event_third")
	_, ok := emitter.listeners.Load("event_first")
	assert.False(t, ok)
	_, ok = emitter.listeners.Load("event_second")
	assert.True(t, ok)
	_, ok = emitter.listeners.Load("event_third")
	assert.False(t, ok)

	// Register events.
	emitter.On("test_1", func() {})
	emitter.AddListener("test_2", func() {})
	emitter.On("test_3", func() {})

	// Remove events.
	emitter.Clear("test_1", "test_2")
	_, ok = emitter.listeners.Load("test_1")
	assert.False(t, ok)
	_, ok = emitter.listeners.Load("test_2")
	assert.False(t, ok)
	_, ok = emitter.listeners.Load("test_3")
	assert.True(t, ok)
	_, ok = emitter.listeners.Load("event_second")
	assert.True(t, ok)

	// Remove all events.
	emitter.RemoveAllListeners()
	_, ok = emitter.listeners.Load("event_second")
	assert.False(t, ok)
	_, ok = emitter.listeners.Load("test_3")
	assert.False(t, ok)

	// Register events.
	emitter.On("test_clear_1", func() {})
	emitter.AddListener("test_clear_2", func() {})
	emitter.On("test_clear_3", func() {})
	_, ok = emitter.listeners.Load("test_clear_1")
	assert.True(t, ok)
	_, ok = emitter.listeners.Load("test_clear_2")
	assert.True(t, ok)
	_, ok = emitter.listeners.Load("test_clear_3")
	assert.True(t, ok)

	// Remove all events.
	emitter.Clear()
	_, ok = emitter.listeners.Load("test_clear_1")
	assert.False(t, ok)
	_, ok = emitter.listeners.Load("test_clear_2")
	assert.False(t, ok)
	_, ok = emitter.listeners.Load("test_clear_3")
	assert.False(t, ok)
}

func TestEmit(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(7)

	emitter := New()

	// Emitting an event without arguments.
	emitter.On("empty", func() {
		defer wg.Done()
	})
	assert.NoError(t, emitter.Emit("empty"))

	// Emitting an event with arguments.
	emitter.On("arguments", func(a, b int) {
		defer wg.Done()

		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
	})
	assert.NoError(t, emitter.Emit("arguments", 1, 2))

	// Emitting an event with multiple arguments.
	emitter.On("multiple", func(i int, s []string, b bool) {
		defer wg.Done()

		assert.Equal(t, 42, i)
		assert.Equal(t, []string{"t", "e", "s", "t"}, s)
		assert.Equal(t, true, b)
	})
	assert.NoError(t, emitter.Emit("multiple", 42, []string{"t", "e", "s", "t"}, true))

	// Emitting an event with multiple arguments.
	emitter.On("variadic_any", func(args ...any) {
		defer wg.Done()

		assert.Equal(t, 42, args[0])
		assert.Equal(t, []string{"t", "e", "s", "t"}, args[1])
		assert.Equal(t, true, args[2])
	})
	assert.NoError(t, emitter.Emit("variadic_any", 42, []string{"t", "e", "s", "t"}, true))

	// Emitting an event with multiple arguments.
	emitter.On("variadic_string", func(s ...string) {
		defer wg.Done()

		assert.Equal(t, "t", s[0])
		assert.Equal(t, "e", s[1])
		assert.Equal(t, "s", s[2])
		assert.Equal(t, "t", s[3])
	})
	assert.NoError(t, emitter.Emit("variadic_string", "t", "e", "s", "t"))

	// Emit event with multiple listeners.
	event_1 := func(a, b int) {
		defer wg.Done()

		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
	}
	event_2 := func(c, d int) {
		defer wg.Done()

		assert.Equal(t, 1, c)
		assert.Equal(t, 2, d)
	}
	emitter.On("multiple_listeners", event_1)
	emitter.On("multiple_listeners", event_2)
	emitter.Emit("multiple_listeners", 1, 2)

	// Emit without event name.
	err := emitter.Emit("")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
	}

	// Emitting an event that doesn't exist.
	err = emitter.Emit("event")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventNotExists, err)
	}

	wg.Wait()
}

func TestEmitSync(t *testing.T) {
	emitter := New()

	// Emitting an event without arguments.
	emitter.On("empty", func() {})
	assert.NoError(t, emitter.EmitSync("empty"))

	// Emitting an event with arguments.
	emitter.On("arguments", func(a, b int) {
		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
	})
	assert.NoError(t, emitter.EmitSync("arguments", 1, 2))

	// Emitting an event with multiple arguments.
	emitter.On("multiple", func(i int, s []string, b bool) {
		assert.Equal(t, 42, i)
		assert.Equal(t, []string{"t", "e", "s", "t"}, s)
		assert.Equal(t, true, b)
	})
	assert.NoError(t, emitter.EmitSync("multiple", 42, []string{"t", "e", "s", "t"}, true))

	// Emitting an event with multiple arguments.
	emitter.On("variadic_any", func(args ...any) {
		assert.Equal(t, 42, args[0])
		assert.Equal(t, []string{"t", "e", "s", "t"}, args[1])
		assert.Equal(t, true, args[2])
	})
	assert.NoError(t, emitter.EmitSync("variadic_any", 42, []string{"t", "e", "s", "t"}, true))

	// Emitting an event.
	emitter.On("any_slice", func(args []any) {
		assert.Equal(t, 42, args[0])
		assert.Equal(t, []string{"t", "e", "s", "t"}, args[1])
		assert.Equal(t, true, args[2])
	})
	assert.NoError(t, emitter.EmitSync("any_slice", []any{42, []string{"t", "e", "s", "t"}, true}))

	// Emitting an event.
	emitter.On("any", func(arg any) {
		assert.Contains(t, []any{5000, "test", testType{}}, arg)
	})
	assert.NoError(t, emitter.EmitSync("any", 5000))
	assert.NoError(t, emitter.EmitSync("any", "test"))
	assert.NoError(t, emitter.EmitSync("any", testType{}))

	// Emitting an event with multiple arguments.
	emitter.On("variadic_string", func(s ...string) {
		assert.Equal(t, "t", s[0])
		assert.Equal(t, "e", s[1])
		assert.Equal(t, "s", s[2])
		assert.Equal(t, "t", s[3])
	})
	assert.NoError(t, emitter.EmitSync("variadic_string", "t", "e", "s", "t"))

	// Emit event with multiple listeners.
	countEvents := 0
	event_1 := func(a, b int) {
		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
		countEvents++
	}
	event_2 := func(c, d int) {
		assert.Equal(t, 1, c)
		assert.Equal(t, 2, d)
		countEvents++
	}
	emitter.On("multiple_listeners", event_1)
	emitter.On("multiple_listeners", event_2)
	emitter.EmitSync("multiple_listeners", 1, 2)
	assert.Equal(t, 2, countEvents)

	// Emit without event name.
	err := emitter.EmitSync("")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
	}

	// Emitting an event that doesn't exist.
	err = emitter.EmitSync("event")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventNotExists, err)
	}
}

func TestPanic(t *testing.T) {
	emitter := New()

	// Wrong number of arguments.
	emitter.On("panic", func(a, b int) {})
	assert.PanicsWithError(t, "Wrong number of arguments. Event panic expected 2 arguments, got 0.", func() {
		emitter.EmitSync("panic")
	})
	assert.PanicsWithError(t, "Wrong number of arguments. Event panic expected 2 arguments, got 1.", func() {
		emitter.EmitSync("panic", 10)
	})
	assert.PanicsWithError(t, "Wrong number of arguments. Event panic expected 2 arguments, got 3.", func() {
		emitter.EmitSync("panic", 10, 20, 30)
	})
	assert.NotPanics(t, func() { emitter.EmitSync("panic", 10, 20) })

	// Wrong number of arguments.
	emitter.On("panic_variadic", func(a int, b ...int) {})
	assert.PanicsWithError(t, "Not enough arguments. Event panic_variadic expected at least 1 arguments, got 0.", func() {
		emitter.EmitSync("panic_variadic")
	})
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic", 10) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic", 10, 20, 30) })

	// Wrong number of arguments.
	emitter.On("panic_variadic_2", func(a int, b string, c ...int) {})
	assert.PanicsWithError(t, "Not enough arguments. Event panic_variadic_2 expected at least 2 arguments, got 0.", func() {
		emitter.EmitSync("panic_variadic_2")
	})
	assert.PanicsWithError(t, "Not enough arguments. Event panic_variadic_2 expected at least 2 arguments, got 1.", func() {
		emitter.EmitSync("panic_variadic_2", 10)
	})
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_2", 10, "test") })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_2", 10, "test", 30) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_2", 10, "test", 30, 40, 50, 60) })

	// Wrong type of arguments.
	emitter.On("panic_variadic_3", func(a int, b ...int) {})
	assert.PanicsWithError(t, "Wrong argument type. Event panic_variadic_3 expected argument 2 to be int, got string.", func() {
		emitter.EmitSync("panic_variadic_3", 10, "test")
	})
	assert.PanicsWithError(t, "Wrong argument type. Event panic_variadic_3 expected argument 2 to be int, got bool.", func() {
		emitter.EmitSync("panic_variadic_3", 10, 20, 30, 40, 50, false, 60)
	})
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_3", 10, 20, 30, 40, 50) })

	// Wrong type of arguments.
	emitter.On("panic_type", func(a int, b string, c testType) {})
	assert.PanicsWithError(t, "Wrong argument type. Event panic_type expected argument 1 to be int, got string.", func() {
		emitter.EmitSync("panic_type", "test", 20, testType{})
	})
	assert.PanicsWithError(t, "Wrong argument type. Event panic_type expected argument 2 to be string, got int.", func() {
		emitter.EmitSync("panic_type", 10, 20, testType{})
	})
	assert.PanicsWithError(t, "Wrong argument type. Event panic_type expected argument 3 to be eventemitter.testType, got bool.", func() {
		emitter.EmitSync("panic_type", 10, "test", false)
	})
	assert.NotPanics(t, func() { emitter.EmitSync("panic_type", 10, "test", testType{}) })

	// Wrong type of arguments.
	emitter.On("panic_any", func(a int, b any) {})
	assert.NotPanics(t, func() { emitter.EmitSync("panic_any", 10, "test") })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_any", 10, true) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_any", 10, testType{}) })

	emitter.On("panic_any_2", func(a bool, b any, c chan int) {})
	assert.PanicsWithError(t, "Wrong argument type. Event panic_any_2 expected argument 3 to be chan int, got chan string.", func() {
		emitter.EmitSync("panic_any_2", false, 100, make(chan string))
	})
	assert.NotPanics(t, func() { emitter.EmitSync("panic_any_2", true, "test", make(chan int)) })

	// Wrong type of arguments.
	emitter.On("panic_variadic_any", func(a string, b bool, c ...any) {})
	assert.PanicsWithError(t, "Wrong argument type. Event panic_variadic_any expected argument 1 to be string, got int.", func() {
		emitter.EmitSync("panic_variadic_any", 10, true)
	})
	assert.PanicsWithError(t, "Wrong argument type. Event panic_variadic_any expected argument 2 to be bool, got string.", func() {
		emitter.EmitSync("panic_variadic_any", "test", "test")
	})
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_any", "test", true) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_any", "test", true, 10, 20, 30) })
	assert.NotPanics(t, func() {
		emitter.EmitSync("panic_variadic_any", "test", true, false, "test", make(chan bool), testType{})
	})
}

func TestEventNames(t *testing.T) {
	emitter := New()

	event_1 := func() {}
	event_2 := func() {}
	event_3 := func() {}
	event_3_1 := func() {}

	emitter.On("event_1", event_1)
	emitter.On("event_2", event_2)
	emitter.On("event_3", event_3)
	emitter.On("event_3", event_3_1)

	names := emitter.EventNames()
	assert.Equal(t, 3, len(names))
	assert.Contains(t, names, "event_1")
	assert.Contains(t, names, "event_2")
	assert.Contains(t, names, "event_3")

	emitter.Off("event_1", event_1)
	emitter.Off("event_2", event_2)
	names = emitter.EventNames()
	assert.Equal(t, 1, len(names))
	assert.Contains(t, names, "event_3")

	emitter.Off("event_3", event_3)
	names = emitter.EventNames()
	assert.Equal(t, 1, len(names))

	emitter.RemoveListener("event_3", event_3_1)
	names = emitter.EventNames()
	assert.Equal(t, 0, len(names))
}

func TestListeners(t *testing.T) {
	emitter := New()

	event_1 := func() {}
	event_2 := func() {}
	event_2_1 := func() {}

	emitter.On("event_1", event_1)
	emitter.On("event_2", event_2)
	emitter.On("event_2", event_2_1)

	// Get listeners for an event.
	listeners, err := emitter.Listeners("event_1")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(listeners))
		assert.NotEmpty(t, listeners)
	}

	listeners, err = emitter.Listeners("event_2")
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(listeners))
		assert.NotEmpty(t, listeners)
	}

	// Get listeners for an event that doesn't exist.
	emitter.RemoveListener("event_1", event_1)
	listeners, err = emitter.Listeners("event_1")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventNotExists, err)
		assert.Nil(t, listeners)
	}

	// Get listeners for an event.
	emitter.RemoveListener("event_2", event_2)
	listeners, err = emitter.Listeners("event_2")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(listeners))
		assert.NotEmpty(t, listeners)
	}

	// Empty event name.
	listeners, err = emitter.Listeners("")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
		assert.Nil(t, listeners)
	}
}

func TestListenerCount(t *testing.T) {
	emitter := New()

	event_1 := func() {}
	event_2 := func() {}
	event_2_1 := func() {}

	emitter.On("event_1", event_1)
	emitter.On("event_2", event_2)
	emitter.On("event_2", event_2_1)

	// Get listeners for an event.
	count, err := emitter.ListenerCount("event_1")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, count)
	}

	count, err = emitter.ListenerCount("event_2")
	if assert.NoError(t, err) {
		assert.Equal(t, 2, count)
	}

	// Get listeners for an event that doesn't exist.
	emitter.RemoveListener("event_1", event_1)
	count, err = emitter.ListenerCount("event_1")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventNotExists, err)
		assert.Equal(t, 0, count)
	}

	// Get listeners for an event.
	emitter.RemoveListener("event_2", event_2)
	count, err = emitter.ListenerCount("event_2")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, count)
	}

	// Empty event name.
	count, err = emitter.ListenerCount("")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
		assert.Equal(t, 0, count)
	}
}
