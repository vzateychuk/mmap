package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func readProcess(amount int) {
	data, file, err := mmapFile()
	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(data)
	defer file.Close()

	lastMsg := ""
	for i := 0; i < amount; i++ {
		msg := string(data[:32])
		if msg != lastMsg {
			fmt.Println("Reader got:", msg)
			lastMsg = msg
		}
		time.Sleep(1 * time.Second)
	}
}

func writeProcess(amount int) {
	data, file, err := mmapFile()
	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(data)
	defer file.Close()

	for i := 0; i < amount; i++ {
		msg := fmt.Sprintf("message %d", i)
		copy(data, msg)
		fmt.Println("Writer sent:", msg)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		panic("usage: go run main.go [write|read]")
	}

	switch os.Args[1] {
	case "write":
		writeProcess(106)
	case "read":
		readProcess(150)
	default:
		panic("unknown command: " + os.Args[1])
	}
}
