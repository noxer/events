package events

import "reflect"

type chanHandler struct {
	typ reflect.Type
	ch  reflect.Value
}

func (h chanHandler) HandleEvent(typ reflect.Type, val reflect.Value) {
	if typ.ConvertibleTo(typ) {
		h.ch.TrySend(val)
	}
}
