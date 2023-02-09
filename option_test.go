package option

import (
	"fmt"
	"io/fs"
	"os"
	"testing"
)

func TestIsNil(t *testing.T) {
	i := Some(30)
	j := None[int]()

	if i.IsNil() {
		t.Error("int '30' can not be nil")
	}

	if !j.IsNil() {
		t.Error("nil int can not have a value")
	}
}

func TestPrimitiveIsEqual(t *testing.T) {
	i := Some(30)
	j := Some(30)

	if i.Unwrap() != j.Unwrap() {
		t.Error("30 != 30")
	}
}

func TestOr(t *testing.T) {
	nilValue := None[*int]()
	twenty := 20
	defaultValue := nilValue.Or(&twenty)

	if *defaultValue != 20 {
		t.Error("Could not get 20")
	}
}

func TestUnwrapOrElse(t *testing.T) {
	f := SomePair(os.ReadFile("./test_data/file_does_not_exist.txt"))

	byt := f.UnwrapOrElse(func(err error) []byte {
		fmt.Println(err.Error())
		return []byte{}
	})

	if len(byt) > 0 {
		t.Error("Should not have reached here")
	}

	fmt.Println(byt)
}

func TestUnwrapOrElseInterfaces(t *testing.T) {
	outDir := "out"
	dirInfo := SomePair(os.Stat(outDir))

	defer func() {
		if err := recover(); err != nil {
			t.Error("UnwrapOrElse with interfaces should work")
		}
	}()

	path := dirInfo.UnwrapOrElse(func(err error) fs.FileInfo {
		if os.IsNotExist(err) {
			crDir := Some(os.Mkdir(outDir, 0755))
			crDir.ExpectNil("could not create directory '" + outDir + "'")
		}

		opt := SomePair(os.Stat(outDir))

		return opt.Expect("unable to create output directory")
	})

	fmt.Println(path.Name())
}

func TestUnwrapNil(t *testing.T) {
	n := None[int]()

	defer func() {
		if err := recover(); err == nil {
			t.Error("Unwrap nil must not work")
		}
	}()

	n.Unwrap()
}

func TestExpect(t *testing.T) {
	e := None[int]()

	defer func() {
		if err := recover(); err == nil {
			t.Error("Expect nil must not work")
		}
	}()

	e.Expect("Int is nil")
}

func TestRealWorld(t *testing.T) {
	f := SomePair(os.ReadFile("./test_data/file.txt"))

	if string(f.Expect("Could not read file")) != "Hello Go Options" {
		t.Error("Failed to read file")
	}

	f = SomePair(os.ReadFile("./test_data/file_does_not_exist.txt"))

	if f.Error() == "" {
		t.Error("Reading non-existent file should have errored")
	}
}
