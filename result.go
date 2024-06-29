package result

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// Result is a type that represents the output of a function in a struct that handles error, value and nil values returned
type Result[T any] struct {
	data *T
	err  error
}

// Expect is similar to Rust's .expect() where the program panics with the given error string if the Result contains
// no data
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

// ExpectNil is an experimental function that only panics when the Result is not a None[T]
//
// This is only ever useful in cases where a None[T] is required.
func (o *Result[T]) ExpectNil(err string) {
	if o.data != nil {
		panic(err)
	}
}

// Or is a function that returns an alternative value if the Result is a None[T] or an Err[T]
//
// This is useful for cases where a default value suffices in the case of a failed operation.
func (o *Result[T]) Or(value T) T {
	if o.data == nil {
		return value
	}

	return *o.data
}

// UnwrapOrElse is a function that allows the mapping of `Err[T]`s to `T`s
//
// This is useful in cases where a function can return multiple types of errors
// and different operations need to be performed in each error case
func (o *Result[T]) UnwrapOrElse(errFn func(err error) T) T {
	if o.err != nil {
		return errFn(o.err)
	}

	return *o.data
}

// Unwrap is similar to Rust's .unwrap() where the program panics if the Result is any of Err[T] or None[T]
//
// In the case of an Err[T] the provided error message is the panic message
//
// In the case of a None[T] `result value is nil` is the panic message
func (o *Result[T]) Unwrap() T {
	if o.data == nil && o.err == nil {
		panic("result value is nil")
	} else if o.data == nil {
		panic(o.err)
	}

	return *o.data
}

// IsNil is a utility function for checking whether a Result is an Err[T] or None[T]
func (o *Result[T]) IsNil() bool {
	return o.data == (*T)(nil)
}

// Error returns the error message contained within an Err[T]
//
// NOTE: Error returns an empty string in the case of a None[T]
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

func innerTry[T any](fn func() T, errorOut *string) *Result[T] {
	defer func() {
		if err := recover(); err != nil {
			*errorOut = fmt.Sprintf("%v", err)
		}
	}()

	rs := fn()

	return &Result[T]{
		data: &rs,
		err:  nil,
	}
}

// None is the unit primitive denoting an empty value
//
// Unwrapping a `None` results in an error
func None[T any]() *Result[T] {
	return &Result[T]{
		data: nil,
		err:  errors.New("value is None"),
	}
}

// Err denotes an error result, `Unwrap`ing or `Expect`ing an Err causes a panic
func Err[T any](err_ error) *Result[T] {
	return &Result[T]{
		data: nil,
		err:  err_,
	}
}

// Some is a Result containing a pointer to a value of type T
//
// A Some Result can be `Unwrap`ed and `Expect`ed
func Some[T any](value T) *Result[T] {
	o := newOptional[T](value)
	return &o
}

// SomePair is a utility function that converts the Golang `value, err` convention
// into a Result of type T where T is the type of `value`
func SomePair[T any](value T, err error) *Result[T] {
	o := newOptionalPair[T](value, err)
	return &o
}

// Try is a utility function for wrapping functions that can panic and returning an Err Result in that case.
//
// Try returns a Some[T] where the wrapped function successfully executes
func Try[T any](fn func() T) *Result[T] {
	var err string = ""
	r := innerTry(fn, &err)

	if r == nil {
		return Err[T](errors.New(err))
	}

	return r
}
