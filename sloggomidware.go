package main

import (
	"context"
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

func (s SlogGomidWare) Before() func(...any) []any {
	return func(in ...any) []any {
		if len(in) > 0 {
			for _, i := range in {
				ctx, ok := i.(context.Context)
				if ok {
					s.logger.Info("Found context", "context msg", ctx.Value("message"))
				}

				s.logger.Info("Request", "arg", i)
			}
		}

		return in
	}
}

func (s SlogGomidWare) After() func(...any) []any {
	return func(out ...any) []any {
		if len(out) < 1 {
			return out
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
		return out
	}
}

var _ gomid.GomidWare = (*SlogGomidWare)(nil)
