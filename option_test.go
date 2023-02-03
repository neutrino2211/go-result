package option

import (
	"fmt"
	"os"
	"testing"
)

func TestIsNil(t *testing.T) {
	i := NewOptional[int](30)
	j := NewOptional[int](nil)

	if i.IsNil() {
		t.Error("int '30' can not be nil")
	}

	if !j.IsNil() {
		t.Error("nil int can not have a value")
	}
}

func TestPrimitiveIsEqual(t *testing.T) {
	i := NewOptional[int](30)
	j := NewOptional[int](30)

	if i.Unwrap() != j.Unwrap() {
		t.Error("30 != 30")
	}
}

func TestOr(t *testing.T) {
	nilValue := NewOptional[int](nil)

	defaultValue := nilValue.Or(20)

	if defaultValue != 20 {
		t.Error("Could not get 20")
	}
}

func TestUnwrapOrElse(t *testing.T) {
	f := NewOptionalPair[[]byte](os.ReadFile("./test_data/file_does_not_exist.txt"))

	defer func() {
		if err := recover(); err == nil {
			t.Error("Unwrap nil must not work")
		}
	}()

	byt := f.UnwrapOrElse(func(err error) {
		fmt.Println(err.Error())
		panic("Error")
	})

	fmt.Println(byt)

	t.Error("Should not have reached here")
}

func TestUnwrapNil(t *testing.T) {
	n := NewOptional[int](nil)

	defer func() {
		if err := recover(); err == nil {
			t.Error("Unwrap nil must not work")
		}
	}()

	n.Unwrap()
}

func TestExpect(t *testing.T) {
	e := NewOptional[int](nil)

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expect nil must not work")
		}
	}()

	e.Expect("Int is nil")
}

func TestRealWorld(t *testing.T) {
	f := NewOptionalPair[[]byte](os.ReadFile("./test_data/file.txt"))

	if string(f.Expect("Could not read file")) != "Hello Go Options" {
		t.Error("Failed to read file")
	}

	f = NewOptionalPair[[]byte](os.ReadFile("./test_data/file_does_not_exist.txt"))

	if f.Error() == "" {
		t.Error("Reading non-existent file should have errored")
	}
}
