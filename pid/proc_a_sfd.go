package main

import (
	"os"
	"fmt"
	"unsafe"
	"golang.org/x/sys/unix"
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
	var sigmask unix.Sigset_t

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

	fd, err := unix.Signalfd(-1, &sigmask, unix.SFD_CLOEXEC|unix.SFD_NONBLOCK)
	if err != nil {
		fmt.Println("failed to create a signalfd: %v", err)
		os.Exit(1)
	}

	signalfd := os.NewFile(uintptr(fd), "signalfd")
	defer signalfd.Close()

	sigmask.Val[0] |= 1 << uint(unix.SIGUSR1-1)
	_, err = unix.Signalfd(fd, &sigmask, 0)
	if err != nil {
		fmt.Println("failed to update sigmask on signalfd: %v", err)
		os.Exit(1)
	}
/*
	.Printf("Sending SIGUSR1 in 2 sec")
	t := time.AfterFunc(2 * time.Second, func() {
		unix.Kill(unix.Getpid(), unix.SIGUSR1)
	})
	defer t.Stop()
*/
	fmt.Println("Waiting for SIGUSR1 from signalfd")
	buf := make([]byte, unsafe.Sizeof(unix.SignalfdSiginfo{}))
	n, err := signalfd.Read(buf)
	if err != nil {
		fmt.Println("failed reading siginfo from signalfd: %v", err)
		os.Exit(1)
	}
	info := *(*unix.SignalfdSiginfo)(unsafe.Pointer(&buf[0]))
	fmt.Printf("n: %d, %+v", n, info)
}
