// file: mmap_write_example.go
package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/mmap"
)

func mmap_write_read_example() {
	const fileName = "example.dat"
	const size = 4096

	// Create file if not exists
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Ensure file is large enough
	if err := file.Truncate(size); err != nil {
		panic(err)
	}

	// Memory map the file for reading (exp/mmap only supports read-only mapping)
	reader, err := mmap.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// Write to file using standard file I/O
	copyBuf := []byte("Hello from mmap!")
	_, err = file.WriteAt(copyBuf, 0)
	if err != nil {
		panic(err)
	}

	// Read from memory-mapped region using mmap
	readBuf := make([]byte, 16)
	n, err := reader.ReadAt(readBuf, 0)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	fmt.Println(string(readBuf[:n]))
}
