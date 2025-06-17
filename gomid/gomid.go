package gomid

import (
	"reflect"
)

type Gomid struct {
	// function
	handler  any
	midwares []GomidWare
}

type GomidOption func(*Gomid)

func New(h any, opt ...GomidOption) any {
	gomid := &Gomid{
		handler:  h,
		midwares: make([]GomidWare, 0),
	}
	for _, applyOpt := range opt {
		applyOpt(gomid)
	}

	hType := reflect.TypeOf(h)
	if hType.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	var beforeMids []any
	for _, mid := range gomid.midwares {

		beforeMids = append(beforeMids, mid.Before())
	}
	var afterMids []any
	for _, mid := range gomid.midwares {
		afterMids = append(afterMids, mid.After())
	}

	handlerValue := reflect.ValueOf(h)
	handlerType := reflect.TypeOf(h)

	makeFunc := reflect.MakeFunc(handlerType, func(args []reflect.Value) []reflect.Value {
		// call before middlewares
		for _, mid := range beforeMids {

			reflect.ValueOf(mid).Call(args)
		}

		// call function itself
		r := handlerValue.Call(args)

		// call after middlewares
		for _, mid := range afterMids {
			reflect.ValueOf(mid).Call(r)
		}

		return r
	})

	return makeFunc.Interface()
}

func AddMid(m GomidWare) GomidOption {
	return func(g *Gomid) {
		g.midwares = append(g.midwares, m)
	}
}

type GomidWare interface {
	Before() func(...any)
	After() func(...any)
}
