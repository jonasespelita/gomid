package main

import (
	"github.com/jonasespelita/gomid/gomid"
	"log/slog"
	"os"
	"reflect"
)

type SlogGomidWare struct {
	logger *slog.Logger
}

func NewSlogGomidWare() gomid.GomidWare {
	return SlogGomidWare{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (s SlogGomidWare) Before() func(...any) {
	return func(in ...any) {
		s.logger.Info("Request", "args", in)
	}
}

func (s SlogGomidWare) After() func(...any) {
	return func(out ...any) {
		if len(out) < 1 {
			return
		}

		r := out[0]

		// if r is Ptr dereference first
		if reflect.TypeOf(r).Kind() == reflect.Ptr {
			r = reflect.ValueOf(r).Elem().Interface()
		}

		s.logger.Info("Response", "return", r)

		if len(out) > 1 {
			// assume second arg is error
			if out[1] != nil {
				err := out[1].(error)
				s.logger.Error("got error", err)
			}

		}

	}
}

var _ gomid.GomidWare = (*SlogGomidWare)(nil)
