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

const devShm = "/dev/shm/"


func main () {
	var ba []byte
	var offset  int64
	var pid_buf [4]byte

	pid := os.Getpid()
	fmt.Println("proc a pid: ", pid)

	flag:= os.O_RDONLY
	pidfil, err := os.OpenFile("proc_a.pid", flag, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		os.Exit(1)
	}


//	ba = cvtInt32toB(uint32(pid))
	ba = pid_buf[:]
	_, err = pidfil.Read(ba)
	if err != nil {
		fmt.Println("error reading file: ", err)
		os.Exit(1)
	}

	pid_a := cvtBtoInt32(pid_buf[:4])
	fmt.Println("proc_a pid: ", pid_a)

	pidfil.Close()
	fmt.Println("wrote pid into pod file!")

// shared file

	shfilnam := "/dev/shm/prr"

	flag = os.O_RDONLY

	shfil, err := os.OpenFile(shfilnam, flag, 0667)
	if err != nil {
		fmt.Println("shared file open error: ", err)
		os.Exit(1)
	}
	fmt.Println("shared file: ", shfil)

	shfd := int(shfil.Fd())

/*
	shnam := "testprr"
	shflag := unix.MFD_CLOEXEC
	shfd , err := unix.MemfdCreate(shnam, shflag)
	if err != nil {
		fmt.Println("MemfdCreate error!", err)
		os.Exit(1)
	}


	err = unix.Ftruncate(shfd, 256)
	if err != nil {
		fmt.Println("ftruncate error!", err)
		os.Exit(1)
	}
*/

	offset = 0
	shlen := 256
	prot := unix.PROT_READ
//| unix.PROT_WRITE
	shflags:= unix.MAP_SHARED

	ba,err  = unix.Mmap(shfd, offset, shlen, prot, shflags)
	if err != nil {
		fmt.Println("Mmap error!", err)
		os.Exit(1)
	}
	fmt.Println("memory map established!")

	fmt.Println("ba proc_b: ", string(ba))

    fmt.Println("proc_a pid: ", pid_a)
    fmt.Println("sending signal to proc_a!", pid_a)
    err = unix.Kill(int(pid_a), unix.SIGUSR1)
    if err != nil {
        fmt.Println("error sending signal: ", err)
        os.Exit(1)
    }
    fmt.Println("signal usr1 sent!")

/*
	fmt.Println("closing shared area!")
	err = os.Remove(shfilnam)
	if err != nil {
		fmt.Println("unlink error!", err)
		os.Exit(1)
	}
*/
}
