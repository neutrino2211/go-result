package main

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/neutrino2211/go-result"
)

func ListDirectory(path string) {
	entries := result.SomePair(os.ReadDir(path)).Or([]fs.DirEntry{})

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("DIR\t%s/%s\n", path, entry.Name())
			ListDirectory(fmt.Sprintf("%s/%s", path, entry.Name()))
		} else {
			fmt.Printf("FILE\t%s/%s\n", path, entry.Name())
		}
	}
}

func main() {
	ListDirectory("../..")
}
