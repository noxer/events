package events

import (
	"testing"
)

// Tests the direct call to Register
func TestDirectRegister(t *testing.T) {
	Register(func() {})
	if events == nil {
		t.Fatal("Calling Register() did not initialize direct calls event queue!")
	}
}

// Tests the direct call to Unregister
func TestDirectUnregister(t *testing.T) {
	Unregister(nil)
	if events == nil {
		t.Fatal("Calling Unregister() did not initialize direct calls event queue!")
	}
}

// Tests the direct call to Event
func TestDirectEvent(t *testing.T) {
	Event(nil)
	if events == nil {
		t.Fatal("Calling Event() did not initialize direct calls event queue!")
	}
}

////////////////////////////////////////////////////////////////////////////////

func TestCreateEventQueue(t *testing.T) {
	q := CreateEventQueue()
	if q == nil {
		t.Error("CreateEventQueue() returned nil!")
	}
	if q.handlers == nil {
		t.Fatal("The handlers list is nil!")
	}
}

type testHandler struct {
	t *testing.T
	d bool
}

func (t *testHandler) HandleEvent(i interface{}) {
	if s, ok := i.(string); ok {
		if s != "testing" {
			t.t.Fatal("The returned string was not \"testing\"!")
		}
	} else {
		t.t.Fatal("The returned interface was not a string!")
	}

	t.d = true
}

func TestEventHandler(t *testing.T) {
	q := CreateEventQueue()
	h := &testHandler{t: t, d: false}

	q.Register(h)

	if q.handlers.Len() != 1 {
		t.Fatal("The event hander count is not 1!")
	}

	q.Event("testing")

	if !h.d {
		t.Fatal("The handler was not called!")
	}

	q.Unregister(h)
	if q.handlers.Len() != 0 {
		t.Fatal("The event handler count is not 0!")
	}
}
