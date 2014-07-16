package events

import "reflect"

type funcHandler struct {
	typ reflect.Type
	fn  reflect.Value
}

func (h funcHandler) HandleEvent(typ reflect.Type, val reflect.Value) {
	if typ == nil {
		h.fn.Call([]reflect.Value{})
	} else if typ.ConvertibleTo(h.typ) {
		h.fn.Call([]reflect.Value{val})
	}
}
