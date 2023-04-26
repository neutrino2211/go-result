package option

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var errYouIdiot error = errors.New("you forgot to panic you idiot")

type Optional[T any] struct {
	data *T
	err  error
}

func (o *Optional[T]) Expect(err string) T {
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

func (o *Optional[T]) ExpectNil(err string) {
	if o.data != nil {
		panic(err)
	}
}

func (o *Optional[T]) Or(value T) T {
	if o.data == nil {
		return value
	}

	return *o.data
}

func (o *Optional[T]) UnwrapOrElse(errFn func(err error) T) T {
	if o.err != nil {
		return errFn(o.err)
	}

	return *o.data
}

func (o *Optional[T]) Unwrap() T {
	if o.data == nil {
		panic("optional value is nil")
	}

	return *o.data
}

func (o *Optional[T]) IsNil() bool {
	return o.data == (*T)(nil)
}

func (o *Optional[T]) Error() string {
	if o.err != nil {
		return o.err.Error()
	}

	return ""
}

func newOptional[T any](value interface{}) Optional[T] {
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

		return Optional[T]{
			data: &cast,
		}
	}

	if value == nil {
		return Optional[T]{
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

	return Optional[T]{
		data: &cast,
	}
}

func newOptionalPair[T any](value interface{}, err error) Optional[T] {
	if err != nil {
		return Optional[T]{
			data: nil,
			err:  err,
		}
	}

	return newOptional[T](value)
}

func None[T any]() *Optional[T] {
	return &Optional[T]{
		data: nil,
		err:  errors.New("value is None"),
	}
}

func Err[T any](err_ error) *Optional[T] {
	return &Optional[T]{
		data: nil,
		err:  err_,
	}
}

func Some[T any](value T) *Optional[T] {
	o := newOptional[T](value)
	return &o
}

func SomePair[T any](value T, err error) *Optional[T] {
	o := newOptionalPair[T](value, err)
	return &o
}
