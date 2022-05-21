package eventemitter

import (
	"fmt"
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

	event := func() {}
	err = emitter.AddListener("other_event", event)
	if assert.NoError(t, err) {
		_, ok := emitter.listeners.Load("other_event")
		assert.True(t, ok)
	}

	err = emitter.On("other_event", func() bool {
		return true
	})
	if assert.NoError(t, err) {
		_, ok := emitter.listeners.Load("event")
		assert.True(t, ok)
	}

	err = emitter.On("other_event", &event)
	if assert.NoError(t, err) {
		listeners, ok := emitter.listeners.Load("other_event")
		assert.Equal(t, 3, len(listeners.([]any)))
		assert.True(t, ok)
	}

	// Register same event multiple times.
	for i := 0; i < 10; i++ {
		err = emitter.AddListener("same_event", event)
		assert.NoError(t, err)
	}
	listeners, ok := emitter.listeners.Load("same_event")
	assert.Equal(t, 10, len(listeners.([]any)))
	assert.True(t, ok)

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

	// Not a function.
	err = emitter.On("nil", "not a function")
	if assert.Error(t, err) {
		assert.Equal(t, ErrNotAFunction, err)
	}

	newEvent := "not a function"
	err = emitter.AddListener("nil", &newEvent)
	if assert.Error(t, err) {
		assert.Equal(t, ErrNotAFunction, err)
	}
}

func TestRemoveListenerOff(t *testing.T) {
	emitter := New()

	event := func() {}
	event_2 := func() {}

	// Register events.
	emitter.AddListener("event_first", event)
	emitter.AddListener("event_second", event)
	emitter.AddListener("event_third", event)
	emitter.AddListener("event_third", &event_2)

	// Remove an event.
	ok, err := emitter.Off("event_first", event)
	if assert.NoError(t, err) {
		assert.True(t, ok)

		_, ok := emitter.listeners.Load("event_first")
		assert.False(t, ok)
	}

	ok, err = emitter.RemoveListener("event_third", event)
	if assert.NoError(t, err) {
		listeners, ok := emitter.listeners.Load("event_third")
		assert.True(t, ok)
		assert.Equal(t, 1, len(listeners.([]any)))
	}

	_, err = emitter.Off("event_third", &event_2)
	if assert.NoError(t, err) {
		_, ok := emitter.listeners.Load("event_third")
		assert.False(t, ok)
	}

	// Event listener does not exist.
	emitter.AddListener("event_missing_listener", event)
	ok, err = emitter.RemoveListener("event_missing_listener", func() {})
	if assert.NoError(t, err) {
		assert.False(t, ok)
	}

	// Event does not exist.
	ok, err = emitter.RemoveListener("event_fourth", func() {})
	if assert.Error(t, err) {
		assert.False(t, ok)
		assert.Equal(t, ErrEventNotExists, err)
	}

	// Empty event name.
	ok, err = emitter.Off("", func() {})
	if assert.Error(t, err) {
		assert.False(t, ok)
		assert.Equal(t, ErrEmptyName, err)
	}

	// Not a function.
	ok, err = emitter.RemoveListener("event_second", "not a function")
	if assert.Error(t, err) {
		assert.False(t, ok)
		assert.Equal(t, ErrNotAFunction, err)
	}

	eventNotAFn := "not a function"
	ok, err = emitter.Off("event_second", &eventNotAFn)
	if assert.Error(t, err) {
		assert.False(t, ok)
		assert.Equal(t, ErrNotAFunction, err)
	}
}

func TestRemoveAllListenersClear(t *testing.T) {
	emitter := New()

	// Register events.
	event := func() {}
	events := []string{"event_1", "event_2", "event_3"}
	for _, eventName := range events {
		emitter.On(eventName, func() {})
		emitter.On(eventName, &event)
	}

	// Remove multiple events.
	emitter.RemoveAllListeners("event_1", "event_2")
	for _, eventName := range []string{"event_1", "event_2"} {
		_, ok := emitter.listeners.Load(eventName)
		assert.False(t, ok)
	}

	listeners, ok := emitter.listeners.Load("event_3")
	assert.True(t, ok)
	assert.Equal(t, 2, len(listeners.([]any)))

	// Register events.
	newEvent := func() {}
	events = []string{"test_1", "test_2", "test_3"}
	for _, eventName := range events {
		emitter.On(eventName, func() {})
		emitter.On(eventName, &newEvent)
	}

	// Remove all events.
	emitter.Clear()
	for _, eventName := range events {
		_, ok := emitter.listeners.Load(eventName)
		assert.False(t, ok)
	}
}

func TestEmit(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(11)

	emitter := New()

	// Emitting an event without arguments.
	emitter.On("empty", func() {
		defer wg.Done()
	})
	assert.NoError(t, emitter.Emit("empty"))

	// Emitting an event with arguments.
	eventArgs := func(a, b int) {
		defer wg.Done()

		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
	}
	emitter.On("arguments", eventArgs)
	emitter.On("arguments", &eventArgs)
	assert.NoError(t, emitter.Emit("arguments", 1, 2))

	// Emitting an event with multiple arguments.
	eventMultipleArgs := func(i int, s []string, b bool) {
		defer wg.Done()

		assert.Equal(t, 42, i)
		assert.Equal(t, []string{"t", "e", "s", "t"}, s)
		assert.Equal(t, true, b)
	}
	emitter.On("multiple", eventMultipleArgs)
	emitter.On("multiple", &eventMultipleArgs)
	assert.NoError(t, emitter.Emit("multiple", 42, []string{"t", "e", "s", "t"}, true))

	// Emitting an event with multiple arguments.
	eventVariadic := func(args ...any) {
		defer wg.Done()

		assert.Equal(t, 42, args[0])
		assert.Equal(t, []string{"t", "e", "s", "t"}, args[1])
		assert.Equal(t, true, args[2])
	}
	emitter.On("variadic_any", eventVariadic)
	emitter.On("variadic_any", &eventVariadic)
	assert.NoError(t, emitter.Emit("variadic_any", 42, []string{"t", "e", "s", "t"}, true))

	// Emitting an event with multiple arguments.
	eventVariadicMultiple := func(a int, b bool, c ...string) {
		defer wg.Done()

		assert.Equal(t, 10, a)
		assert.Equal(t, true, b)
		assert.Equal(t, []string{"t", "e", "s", "t"}, c)
	}
	emitter.On("variadic", eventVariadicMultiple)
	emitter.On("variadic", &eventVariadicMultiple)
	assert.NoError(t, emitter.Emit("variadic", 10, true, "t", "e", "s", "t"))

	// Emit event with multiple listeners.
	event_1 := func(a int, b string) {
		defer wg.Done()

		assert.Equal(t, 1, a)
		assert.Equal(t, "test", b)
	}
	event_2 := func(a int, b string) {
		defer wg.Done()

		assert.Equal(t, 1, a)
		assert.Equal(t, "test", b)
	}
	emitter.On("different_listeners", event_1)
	emitter.On("different_listeners", &event_2)
	emitter.Emit("different_listeners", 1, "test")

	// Emit without event name.
	err := emitter.Emit("")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
	}

	// Emitting an event that doesn't exist.
	err = emitter.Emit("event_not_exists")
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
	eventArgs := func(a, b int) {
		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
	}
	emitter.On("arguments", eventArgs)
	emitter.On("arguments", &eventArgs)
	assert.NoError(t, emitter.EmitSync("arguments", 1, 2))

	// Emitting an event with multiple arguments.
	eventVariadicMultiple := func(a int, b bool, c ...string) {
		assert.Equal(t, 10, a)
		assert.Equal(t, true, b)
		assert.Equal(t, []string{"t", "e", "s", "t"}, c)
	}
	emitter.On("variadic", eventVariadicMultiple)
	emitter.On("variadic", &eventVariadicMultiple)
	assert.NoError(t, emitter.EmitSync("variadic", 10, true, "t", "e", "s", "t"))

	// Emit event with multiple listeners.
	event_1 := func(a int, b string) {
		assert.Equal(t, 1, a)
		assert.Equal(t, "test", b)
	}
	event_2 := func(a int, b string) {
		assert.Equal(t, 1, a)
		assert.Equal(t, "test", b)
	}
	emitter.On("different_listeners", event_1)
	emitter.On("different_listeners", &event_2)
	emitter.EmitSync("different_listeners", 1, "test")

	// Emit without event name.
	err := emitter.EmitSync("")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
	}

	// Emitting an event that doesn't exist.
	err = emitter.EmitSync("event_not_exists")
	if assert.Error(t, err) {
		assert.Equal(t, ErrEventNotExists, err)
	}
}

func createTypeErr(event string, pos int, expected string, got string) string {
	return fmt.Sprintf("Wrong argument type. Event %s expected argument %d to be %s, got %s.", event, pos, expected, got)
}

func TestPanic(t *testing.T) {
	emitter := New()

	// Wrong number of arguments.
	emitter.On("panic", func(a, b int) {})
	assert.PanicsWithError(t, (&argsError{"panic", 2, 0}).Error(), func() { emitter.EmitSync("panic") })
	assert.PanicsWithError(t, (&argsError{"panic", 2, 1}).Error(), func() { emitter.EmitSync("panic", 10) })
	assert.PanicsWithError(t, (&argsError{"panic", 2, 3}).Error(), func() { emitter.EmitSync("panic", 10, 20, 30) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic", 10, 20) })

	// Wrong number of arguments.
	emitter.On("panic_variadic", func(a int, b ...int) {})
	assert.PanicsWithError(t, (&argsError{"panic_variadic", 1, 0}).Error(), func() { emitter.EmitSync("panic_variadic") })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic", 10) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic", 10, 20, 30) })

	// Wrong number of arguments.
	emitter.On("panic_variadic_any", func(a int, b string, c ...any) {})
	assert.PanicsWithError(t, (&argsError{"panic_variadic_any", 2, 0}).Error(), func() { emitter.EmitSync("panic_variadic_any") })
	assert.PanicsWithError(t, (&argsError{"panic_variadic_any", 2, 1}).Error(), func() { emitter.EmitSync("panic_variadic_any", 10) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_any", 10, "test") })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_any", 10, "test", 30) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_any", 10, "test", 30, true, make(chan int), &emitter) })

	// Wrong type of arguments.
	emitter.On("panic_variadic_type", func(a int, b ...int) {})
	assert.PanicsWithError(t, createTypeErr("panic_variadic_type", 2, "int", "string"), func() { emitter.EmitSync("panic_variadic_type", 10, "test") })
	assert.PanicsWithError(t, createTypeErr("panic_variadic_type", 2, "int", "bool"), func() { emitter.EmitSync("panic_variadic_type", 10, 20, 30, 40, 50, false, 60) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_variadic_3", 10, 20, 30, 40, 50) })

	// Wrong type of arguments.
	emitter.On("panic_type", func(a int, b string, c testType) {})
	assert.PanicsWithError(t, createTypeErr("panic_type", 1, "int", "string"), func() { emitter.EmitSync("panic_type", "test", 20, testType{}) })
	assert.PanicsWithError(t, createTypeErr("panic_type", 2, "string", "int"), func() { emitter.EmitSync("panic_type", 10, 20, testType{}) })
	assert.PanicsWithError(t, createTypeErr("panic_type", 3, "eventemitter.testType", "bool"), func() { emitter.EmitSync("panic_type", 10, "test", false) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_type", 10, "test", testType{}) })

	// Wrong type of arguments.
	emitter.On("panic_any", func(a int, b any) {})
	assert.NotPanics(t, func() { emitter.EmitSync("panic_any", 10, "test") })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_any", 10, true) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_any", 10, testType{}) })

	emitter.On("panic_any_type", func(a bool, b any, c chan int) {})
	assert.PanicsWithError(t, createTypeErr("panic_any_type", 3, "chan int", "chan string"), func() { emitter.EmitSync("panic_any_type", false, 100, make(chan string)) })
	assert.NotPanics(t, func() { emitter.EmitSync("panic_any_type", true, "test", make(chan int)) })
}

func TestEventNames(t *testing.T) {
	emitter := New()

	event := func() {}
	events := []string{"event_1", "event_2", "event_3"}
	for _, eventName := range events {
		emitter.On(eventName, event)
	}

	// Get event names.
	names := emitter.EventNames()
	assert.Equal(t, 3, len(names))
	for _, eventName := range events {
		assert.Contains(t, names, eventName)
	}

	emitter.Off("event_1", event)
	emitter.Off("event_2", event)

	names = emitter.EventNames()
	assert.Equal(t, 1, len(names))
	assert.Contains(t, names, "event_3")

	emitter.Off("event_3", event)
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
	emitter.On("event_2", &event_2_1)

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
	emitter.On("event_2", &event_2_1)

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

func TestIsFunction(t *testing.T) {
	emitter := New()

	event_1 := func() {}
	event_2 := &event_1

	assert.True(t, emitter.isFunction(event_1))
	assert.True(t, emitter.isFunction(event_2))

	notAfunction := "not a function"
	for _, v := range []any{1, "test", []*any{}, true, &notAfunction, []any{}, nil, notAfunction, make(chan int), 42.00} {
		assert.False(t, emitter.isFunction(v))
	}
}

func TestIsEqual(t *testing.T) {
	emitter := New()

	event_1 := func() {}
	event_2 := &event_1
	event_3 := event_1
	event_4 := event_2

	event_other_1 := func() {}
	event_other_2 := &event_other_1
	event_other_3 := event_other_1

	assert.True(t, emitter.isEqual(event_1, event_1))
	assert.True(t, emitter.isEqual(event_1, event_3))
	assert.True(t, emitter.isEqual(event_2, event_2))
	assert.True(t, emitter.isEqual(event_2, event_4))
	assert.False(t, emitter.isEqual(event_1, event_2))
	assert.False(t, emitter.isEqual(event_2, event_3))

	assert.False(t, emitter.isEqual(event_1, event_other_1))
	assert.False(t, emitter.isEqual(event_2, event_other_2))
	assert.False(t, emitter.isEqual(event_3, event_other_3))
}

func BenchmarkAddListener(b *testing.B) {
	emitter := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.AddListener("event", func() {})
	}
}

func BenchmarkAddListenerPointer(b *testing.B) {
	emitter := New()
	testEvent := func() {}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.AddListener("event", &testEvent)
	}
}

func BenchmarkEmitSync(b *testing.B) {
	emitter := New()
	emitter.AddListener("event", func() {})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.EmitSync("event")
	}
}

func BenchmarkEmitSyncPointer(b *testing.B) {
	emitter := New()
	testEvent := func() {}
	emitter.AddListener("event", &testEvent)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.EmitSync("event")
	}
}

func BenchmarkEmitSyncWithArguments(b *testing.B) {
	emitter := New()
	emitter.AddListener("event", func(a int, b string, c bool) {})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.EmitSync("event", 1, "test", true)
	}
}

func BenchmarkEmitSyncPointerWithArguments(b *testing.B) {
	emitter := New()
	testEvent := func(a int, b string, c bool) {}
	emitter.AddListener("event", &testEvent)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.EmitSync("event", 1, "test", true)
	}
}

func BenchmarkEmitSyncVariadic(b *testing.B) {
	emitter := New()
	emitter.AddListener("event", func(a int, b bool, c ...any) {})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.EmitSync("event", 1, true, "test", false, 1000)
	}
}

func BenchmarkEmitSyncPointerVariadic(b *testing.B) {
	emitter := New()
	testEvent := func(a int, b bool, c ...any) {}
	emitter.AddListener("event", &testEvent)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.EmitSync("event", 1, true, "test", false, 1000)
	}
}

func BenchmarkRemoveListener(b *testing.B) {
	emitter := New()

	event := func() {}
	for i := 0; i < b.N; i++ {
		emitter.AddListener("event", event)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.RemoveListener("event", event)
	}
}

func BenchmarkRemoveListenerPointer(b *testing.B) {
	emitter := New()
	testEvent := func() {}

	for i := 0; i < b.N; i++ {
		emitter.AddListener("event", &testEvent)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emitter.RemoveListener("event", &testEvent)
	}
}
