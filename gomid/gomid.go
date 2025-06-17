package gomid

import (
	"reflect"
)

type gomid struct {
	// function
	handler  any
	midwares []GomidWare
}

type GomidOption func(*gomid)

func New(h any, opt ...GomidOption) any {
	g := &gomid{
		handler:  h,
		midwares: make([]GomidWare, 0),
	}
	for _, applyOpt := range opt {
		applyOpt(g)
	}

	hType := reflect.TypeOf(h)
	if hType.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	var beforeMids []any
	for _, mid := range g.midwares {
		beforeMids = append(beforeMids, mid.Before())
	}

	var afterMids []any
	for _, mid := range g.midwares {
		// prepend after middle wares for FILO
		afterMids = append([]any{mid.After()}, afterMids...)
	}

	handlerValue := reflect.ValueOf(h)
	handlerType := reflect.TypeOf(h)

	makeFunc := reflect.MakeFunc(handlerType, func(args []reflect.Value) []reflect.Value {

		// call before middlewares
		for _, mid := range beforeMids {
			preCallR := reflect.ValueOf(mid).Call(args)
			// recreate []Value for returning from Before() call
			if len(preCallR) > 0 {
				returns := preCallR[0].Interface().([]any)
				for i, a := range returns {
					if a != nil {
						args[i] = reflect.ValueOf(a)
					} else {
						args[i] = reflect.Zero(args[i].Type())
					}
				}
			}
		}

		// call function itself
		r := handlerValue.Call(args)
		// call after middlewares
		for _, mid := range afterMids {
			postCallR := reflect.ValueOf(mid).Call(r)
			// recreate []Value for returning from After() call
			if len(postCallR) > 0 {
				returns := postCallR[0].Interface().([]any)
				for i, a := range returns {
					if a != nil {
						r[i] = reflect.ValueOf(a)
					} else {
						r[i] = reflect.Zero(r[i].Type())
					}
				}
			}
		}

		return r
	})

	return makeFunc.Interface()
}

func AddMid(m GomidWare) GomidOption {
	return func(g *gomid) {
		g.midwares = append(g.midwares, m)
	}
}

type GomidWare interface {
	Before() func(...any) []any
	After() func(...any) []any
}
