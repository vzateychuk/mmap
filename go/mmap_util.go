package main

import (
	"os"

	"golang.org/x/sys/unix"
)

const (
	fileName = "/tmp/ipc_shared_mem"
	size     = 4096
)

// mmapFile opens or creates a memory-mapped file for inter-process communication.
func mmapFile() ([]byte, *os.File, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, err
	}
	if err := file.Truncate(size); err != nil {
		file.Close()
		return nil, nil, err
	}
	data, err := unix.Mmap(int(file.Fd()), 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		file.Close()
		return nil, nil, err
	}
	return data, file, nil
}
