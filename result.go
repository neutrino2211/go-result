package result

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var errYouIdiot error = errors.New("you forgot to panic you idiot")

type Result[T any] struct {
	data *T
	err  error
}

func (o *Result[T]) Expect(err string) T {
	if o.data == nil {
		panic(err)
	}

	deref := *o.data

	if reflect.ValueOf(deref).Kind() == reflect.String {
		ptr := unsafe.Pointer(&deref)
		castString := (*string)(ptr)

		if *castString == "" {
			panic(err)
		}
	}

	return deref
}

func (o *Result[T]) ExpectNil(err string) {
	if o.data != nil {
		panic(err)
	}
}

func (o *Result[T]) Or(value T) T {
	if o.data == nil {
		return value
	}

	return *o.data
}

func (o *Result[T]) UnwrapOrElse(errFn func(err error) T) T {
	if o.err != nil {
		return errFn(o.err)
	}

	return *o.data
}

func (o *Result[T]) Unwrap() T {
	if o.data == nil && o.err == nil {
		panic("optional value is nil")
	} else if o.data == nil {
		panic(o.err)
	}

	return *o.data
}

func (o *Result[T]) IsNil() bool {
	return o.data == (*T)(nil)
}

func (o *Result[T]) Error() string {
	if o.err != nil {
		return o.err.Error()
	}

	return ""
}

func newOptional[T any](value interface{}) Result[T] {
	var tmp T

	if reflect.ValueOf(value).Kind() == reflect.Ptr {
		cast, ok := value.(T)

		if !ok {
			panic(
				fmt.Sprintf(
					"Failed to create Optional: %s to %s",
					reflect.ValueOf(tmp).Type().Name(),
					reflect.ValueOf(value).Type().Name(),
				),
			)
		}

		return Result[T]{
			data: &cast,
		}
	}

	if value == nil {
		return Result[T]{
			data: nil,
		}
	}

	if reflect.ValueOf(tmp).Kind() != reflect.ValueOf(value).Kind() {
		panic(
			fmt.Sprintf(
				"Failed to create Optional: type mismatch %s to %s",
				reflect.ValueOf(tmp).Type().Name(),
				reflect.ValueOf(value).Type().Name(),
			),
		)
	}

	cast, ok := value.(T)

	if !ok {
		panic("Failed optional cast")
	}

	return Result[T]{
		data: &cast,
	}
}

func newOptionalPair[T any](value interface{}, err error) Result[T] {
	if err != nil {
		return Result[T]{
			data: nil,
			err:  err,
		}
	}

	return newOptional[T](value)
}

func None[T any]() *Result[T] {
	return &Result[T]{
		data: nil,
		err:  errors.New("value is None"),
	}
}

func Err[T any](err_ error) *Result[T] {
	return &Result[T]{
		data: nil,
		err:  err_,
	}
}

func Some[T any](value T) *Result[T] {
	o := newOptional[T](value)
	return &o
}

func SomePair[T any](value T, err error) *Result[T] {
	o := newOptionalPair[T](value, err)
	return &o
}

func innerTry[T any](fn func() T) *Result[T] {
	var ropt error
	defer func() {
		if err := recover(); err != nil {
			ropt = fmt.Errorf("Try Error: #%v", err)
		}
	}()

	rs := fn()

	return &Result[T]{
		data: &rs,
		err:  ropt,
	}
}

func Try[T any](fn func() T) *Result[T] {
	r := innerTry(fn)

	if r == nil {
		return None[T]()
	}

	return r
}
