package main

import (
	"os"
	"fmt"
	"unsafe"
	"golang.org/x/sys/unix"
	"syscall"
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
	var buf [4]byte
	var tstsig syscall.Signal

	pid := os.Getpid()
	fmt.Println("proc_b pid: ", pid)

	flag:= os.O_RDONLY
	pidfil, err := os.OpenFile("proc_a.pid", flag, 0666)
	defer pidfil.Close()
	if err != nil {
		fmt.Println("error opening file: ", err)
		os.Exit(1)
	}

	ba = buf[:]
	_, err = pidfil.Read(ba)
	if err != nil {
		fmt.Println("error reading file: ", err)
		os.Exit(1)
	}

	pid_a := cvtBtoInt32(buf[:4])

	fmt.Println("proc_a pid: ", pid_a)
	fmt.Println("sending signal to proc_a!", pid_a)
//	err = unix.Kill(int(pid_a), unix.SIGUSR1)
	tstsig = 34
	err = unix.Kill(int(pid_a), tstsig)
	if err != nil {
		fmt.Println("error sending signal: ", err)
		os.Exit(1)
	}
	fmt.Println("signal sent!")
}
