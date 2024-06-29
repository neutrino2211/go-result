package main

import (
	"errors"
	"strconv"

	"github.com/neutrino2211/go-result"
)

func PositiveAtoi(str string) *result.Result[int] {
	value, err := strconv.Atoi(str)

	if err != nil {
		return result.Err[int](err) // Return the conversion error
	}

	if value < 0 {
		return result.Err[int](errors.New("PositiveAtoi: number provided must be >= 0")) // Custom failure case
	}

	return result.Some(value)
}

func main() {
	one := PositiveAtoi("1")
	minusOne := PositiveAtoi("-1")
	notANumber := PositiveAtoi("one")

	println(one.Unwrap())
	println(minusOne.Error())
	println(notANumber.Error())
}
