package main

import (
	"os"
	"os/signal"
	"fmt"
	"unsafe"
//	"golang.org/x/sys/unix"
)

func cvtBtoInt32(b []byte) uint32 {
    // equivalnt of return int32(binary.LittleEndian.Uint32(b))
    return uint32(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24)
}

func cvtInt32toB(h uint32)(b []byte) {
    b =(*[4]byte)(unsafe.Pointer(&h))[:]
    return
}

func cvtBtoInt64(b []byte) uint64 {
    return uint64(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
}

func cvtInt64toB(h uint64)(b []byte) {
    b = (*[8]byte)(unsafe.Pointer(&h))[:]
    return
}




func main () {
	var ba []byte

	pid := os.Getpid()
	fmt.Println("proc a pid: ", pid)

	flag:= os.O_RDWR | os.O_CREATE
	pidfil, err := os.OpenFile("proc_a.pid", flag, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		os.Exit(1)
	}

	ba = cvtInt32toB(uint32(pid))

	_, err = pidfil.Write(ba)
	if err != nil {
		fmt.Println("error writing file: ", err)
		os.Exit(1)
	}
	pidfil.Close()
	fmt.Println("wrote pid into pod file!")

	c := make(chan os.Signal, 1)

	// Passing no signals to Notify means that
	// all signals will be sent to the channel.
	signal.Notify(c)

	fmt.Println("waiting for signal!")
	// Block until any signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
}

