package eventemitter

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	emitter.On("same_event", func() {})
	err = emitter.On("same_event", func() {})
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

	// Register an event with a wrong function cal
	err = emitter.On("nil", "not a function")
	if assert.Error(t, err) {
		assert.Equal(t, ErrNotAFunction, err)
	}
}

func TestRemoveListenerOff(t *testing.T) {
	emitter := New()

	// Register events.
	emitter.On("event_first", func() {})
	emitter.AddListener("event_second", func() {})
	emitter.On("event_third", func() {})

	// Remove an event.
	err := emitter.Off("event_first")
	if assert.NoError(t, err) {
		_, ok := emitter.listeners.Load("event_first")
		assert.False(t, ok)
	}

	// Remove an event.
	err = emitter.RemoveListener("event_third")
	if assert.NoError(t, err) {
		_, ok := emitter.listeners.Load("event_third")
		assert.False(t, ok)
	}

	// Event does not exist.
	err = emitter.Off("event_fourth")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventNotExists, err)
	}

	// Empty event name.
	err = emitter.RemoveListener("")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
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
	wg.Add(5)

	emitter := New()

	// Emitting an event without arguments.
	emitter.On("empty", func() {
		defer wg.Done()
	})
	err := emitter.Emit("empty")
	assert.NoError(t, err)

	// Emitting an event with arguments.
	emitter.On("arguments", func(a, b int) {
		defer wg.Done()

		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
	})
	err = emitter.Emit("arguments", 1, 2)
	assert.NoError(t, err)

	// Emitting an event with multiple arguments.
	emitter.On("multiple", func(i int, s []string, b bool) {
		defer wg.Done()

		assert.Equal(t, 42, i)
		assert.Equal(t, []string{"t", "e", "s", "t"}, s)
		assert.Equal(t, true, b)
	})
	err = emitter.Emit("multiple", 42, []string{"t", "e", "s", "t"}, true)
	assert.NoError(t, err)

	// Emitting an event with multiple arguments.
	emitter.On("variadic_any", func(args ...any) {
		defer wg.Done()

		assert.Equal(t, 42, args[0])
		assert.Equal(t, []string{"t", "e", "s", "t"}, args[1])
		assert.Equal(t, true, args[2])
	})
	err = emitter.Emit("variadic_any", 42, []string{"t", "e", "s", "t"}, true)
	assert.NoError(t, err)

	// Emitting an event with multiple arguments.
	emitter.On("variadic_string", func(s ...string) {
		defer wg.Done()

		assert.Equal(t, "t", s[0])
		assert.Equal(t, "e", s[1])
		assert.Equal(t, "s", s[2])
		assert.Equal(t, "t", s[3])
	})
	err = emitter.Emit("variadic_string", "t", "e", "s", "t")
	assert.NoError(t, err)

	// Emit without event name.
	err = emitter.Emit("")
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
	err := emitter.EmitSync("empty")
	assert.NoError(t, err)

	// Emitting an event with arguments.
	emitter.On("arguments", func(a, b int) {
		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
	})
	err = emitter.EmitSync("arguments", 1, 2)
	assert.NoError(t, err)

	// Emitting an event with multiple arguments.
	emitter.On("multiple", func(i int, s []string, b bool) {
		assert.Equal(t, 42, i)
		assert.Equal(t, []string{"t", "e", "s", "t"}, s)
		assert.Equal(t, true, b)
	})
	err = emitter.EmitSync("multiple", 42, []string{"t", "e", "s", "t"}, true)
	assert.NoError(t, err)

	// Emitting an event with multiple arguments.
	emitter.On("variadic_any", func(args ...any) {
		assert.Equal(t, 42, args[0])
		assert.Equal(t, []string{"t", "e", "s", "t"}, args[1])
		assert.Equal(t, true, args[2])
	})
	err = emitter.EmitSync("variadic_any", 42, []string{"t", "e", "s", "t"}, true)
	assert.NoError(t, err)

	// Emitting an event with multiple arguments.
	emitter.On("variadic_string", func(s ...string) {
		assert.Equal(t, "t", s[0])
		assert.Equal(t, "e", s[1])
		assert.Equal(t, "s", s[2])
		assert.Equal(t, "t", s[3])
	})
	err = emitter.EmitSync("variadic_string", "t", "e", "s", "t")
	assert.NoError(t, err)

	// Emit without event name.
	err = emitter.EmitSync("")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
	}

	// Emitting an event that doesn't exist.
	err = emitter.EmitSync("event")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventNotExists, err)
	}
}

func TestEventNames(t *testing.T) {
	emitter := New()

	emitter.On("event_1", func() {})
	emitter.On("event_2", func() {})
	emitter.On("event_3", func() {})

	names := emitter.EventNames()
	assert.Equal(t, 3, len(names))
	assert.Contains(t, names, "event_1")
	assert.Contains(t, names, "event_2")
	assert.Contains(t, names, "event_3")

	emitter.Off("event_1")
	emitter.Off("event_2")
	names = emitter.EventNames()
	assert.Equal(t, 1, len(names))
	assert.Contains(t, names, "event_3")

	emitter.Off("event_3")
	names = emitter.EventNames()
	assert.Equal(t, 0, len(names))
	assert.Empty(t, names)
}
