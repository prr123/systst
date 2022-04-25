package main

import (
	"fmt"
	"os"
//	"io"
	"net"
	"strconv"
	"syscall"
//    "unsafe"
  	"golang.org/x/sys/unix"
)

const SockAddr = "/tmp/echo.sock"

func main() {
   var ba []byte
   var offset  int64

    if err := os.RemoveAll(SockAddr); err != nil {
		fmt.Println("error socket removal!", err)
		os.Exit(1)
    }

    l, err := net.Listen("unix", SockAddr)
    if err != nil {
        fmt.Println("listen error:", err)
		os.Exit(1)
    }

	ulis:= l.(*net.UnixListener)

	ufil, err :=ulis.File()
    if err != nil {
        fmt.Println("file error:", err)
		os.Exit(1)
    }

	fd := ufil.Fd()
	nam := ufil.Name()

	fmt.Println(" fd: ", strconv.Itoa(int(fd)), " Name: ", nam)

    defer l.Close()



//    for {
        // Accept new connections, dispatching them to echoServer
        // in a goroutine.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("accept error:", err)
			os.Exit(1)
        }
		fmt.Println("client connected: ", conn.RemoteAddr().Network())

		rec_msg := make([]byte,128)
		count, err := conn.Read(rec_msg)
        if err != nil {
            fmt.Println("read error:", err)
			os.Exit(1)
        }

		fmt.Println("msg : ", string(rec_msg), " | ", count)
        if err != nil {
            fmt.Println("accept error:", err)
			os.Exit(1)
        }



    ucon, ok := conn.(*net.UnixConn)
    if !ok {
        fmt.Println("unix conn error:")
        os.Exit(1)
    }
    file, err := ucon.File()
    if err != nil {
        fmt.Println("ucon file error:", err)
        os.Exit(1)
    }

    fda:=file.Fd()
    fmt.Println(" fd alt: ", fda)

// need to get the fd of the shred file area
    flag:= os.O_RDWR | os.O_CREATE
    shfilnam := "/dev/shm/prr"

    shfil, err := os.OpenFile(shfilnam, flag, 0667)
    if err != nil {
        fmt.Println("shared file open error: ", err)
        os.Exit(1)
    }
    fmt.Println("shared file: ", shfil)

    shfd := int(shfil.Fd())

//alternative test next
//func MemfdCreate(name string, flags int) (fd int, err error)
//	shflags := unix.MFD_CLOEXEC | unix.MFD_ALLOW_SEALING | unix.MFD_HUGE_...)
/*
	shfd, err := unix.MemfdCreate("shprr", shflags)
    if err != nil {
        fmt.Println("shprr creation error!", err)
        os.Exit(1)
    }

*/
    err = unix.Ftruncate(shfd, 256)
    if err != nil {
        fmt.Println("ftruncate error!", err)
        os.Exit(1)
    }

    filinfo2, err := os.Stat(shfilnam)

    fmt.Println("file info: ", filinfo2.Size())


    offset = 0
    shlen := 256
    prot := unix.PROT_READ | unix.PROT_WRITE
    shflags:= unix.MAP_SHARED
    ba,err  = unix.Mmap(shfd, offset, shlen, prot, shflags)
    if err != nil {
        fmt.Println("Mmap error!", err)
        os.Exit(1)
    }
    fmt.Println("memory map established!")

    copy(ba, "hello shared memory")

//    fmt.Println(" ba: ", ba)
    fmt.Println(" ba: ", string(ba))
/*
    filinfo, err := os.Stat(shfilnam)
    if err != nil {
        fmt.Println("shared file info error: ", err)
        os.Exit(1)
    }
*/
//		var msg_ar [][]byte 
		nmsg := []byte("client: msg received!")
//		msg_ar[0] = nmsg

		rights := syscall.UnixRights(shfd)
    	err = syscall.Sendmsg(int(fda), nmsg, rights, nil, 0)
        if err != nil {
            fmt.Println("sendmsg error:", err)
			os.Exit(1)
        }

		fmt.Println("closing!")
		conn.Close()
//	} // for
}
