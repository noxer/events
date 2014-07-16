package events

import (
	"container/list"
	"reflect"
	"sync"
)

var (
	events *EventQueue
)

func Register(handler interface{}) {
	if events == nil {
		events = CreateEventQueue()
	}
	events.Register(handler)
}

func Unregister(handler interface{}) {
	if events == nil {
		events = CreateEventQueue()
	}
	events.Unregister(handler)
}

func Event(event interface{}) {
	if events == nil {
		events = CreateEventQueue()
	}
	events.Event(event)
}

type internalHandler interface {
	HandleEvent(reflect.Type, reflect.Value)
}

type EventHandler interface {
	HandleEvent(interface{})
}

type EventQueue struct {
	handlers *list.List   // Holds the list of handlers
	mutex    sync.RWMutex // Making all that threadsafe
}

func CreateEventQueue() *EventQueue {
	return &EventQueue{
		handlers: list.New(),
	}
}

func (q *EventQueue) Register(handler interface{}) {
	// Does it implement the handler type
	if evHandler, ok := handler.(EventHandler); ok {
		q.registerHandler(evHandler)
		return
	}

	// Determine the type of "handler"
	typ := reflect.TypeOf(handler)
	val := reflect.ValueOf(handler)
	switch typ.Kind() {
	case reflect.Array, reflect.Slice: // Recursively call this function for every element
		for i := 0; i < val.Len(); i++ {
			q.Register(val.Index(i).Interface())
		}

	case reflect.Chan: // Register a channel
		q.registerChan(typ, val)

	case reflect.Func: // Register a function
		q.registerFunc(typ, val)

	default:
		panic("Trying to register unsupported handler type: \"" + typ.String() + "\"")

	}
}

func (q *EventQueue) Unregister(handler interface{}) {
	if handler == nil {
		return
	}

	// Try to unregister an array
	typ := reflect.TypeOf(handler)
	if typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice {
		val := reflect.ValueOf(handler)
		for i := 0; i < val.Len(); i++ {
			Unregister(val.Index(i).Interface())
		}
		return
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	for e := q.handlers.Front(); e != nil; e = e.Next() {
		switch h := e.Value.(type) {
		case EventHandler: // Remove an event handler
			if h == handler {
				q.handlers.Remove(e)
				return
			}

		case chanHandler: // Remove a channel
			if h.ch.Interface() == handler {
				q.handlers.Remove(e)
				return
			}

		case funcHandler: // Remove a function
			if h.fn.Interface() == handler {
				q.handlers.Remove(e)
				return
			}
		}
	}
}

func (q *EventQueue) Event(event interface{}) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	typ := reflect.TypeOf(event)
	val := reflect.ValueOf(event)

	for e := q.handlers.Front(); e != nil; e = e.Next() {
		switch h := e.Value.(type) {
		case EventHandler:
			h.HandleEvent(event)
		case internalHandler:
			h.HandleEvent(typ, val)
		}
	}
}

func (q *EventQueue) registerHandler(handler EventHandler) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.handlers.PushBack(handler)
}

func (q *EventQueue) registerChan(typ reflect.Type, val reflect.Value) {
	q.mutex.Lock()
	q.mutex.Unlock()

	q.handlers.PushBack(chanHandler{typ: typ.Elem(), ch: val})
}

func (q *EventQueue) registerFunc(typ reflect.Type, val reflect.Value) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if typ.NumIn() == 0 {
		q.handlers.PushBack(funcHandler{fn: val})
	} else if typ.NumIn() == 1 {
		q.handlers.PushBack(funcHandler{typ: typ.In(0), fn: val})
	} else {
		panic("Trying to register a handler function with too many parameters, 0 or 1 permitted")
	}
}
