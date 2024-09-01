package util

import (
	"os"
)

func MustHostname() string {
	return Must(os.Hostname())
}

func Must[T any](t T, err error) T {
	return Result(t, err).Must()
}
func MustDo(err error) {
	Result[any](nil, err).Do()
}

type ErrorHandler[T any] struct {
	Data T
	Err  error
}

func Result[T any](data T, err error) *ErrorHandler[T] {
	return &ErrorHandler[T]{Data: data, Err: err}
}

func (h *ErrorHandler[T]) Must() T {
	if h.Err != nil {
		panic(h.Err)
	}
	return h.Data
}
func (h *ErrorHandler[T]) Do() {
	if h.Err != nil {
		panic(h.Err)
	}
}

func (h *ErrorHandler[T]) RunOk(f func(T)) error {
	if h.Err != nil {
		return h.Err
	}
	f(h.Data)
	return nil
}

func (h *ErrorHandler[T]) Run(f func(T) error) error {
	if h.Err != nil {
		return h.Err
	}
	return f(h.Data)
}
