package main

import (
	"fmt"
	"os"
//	"io"
	"net"
//	"reflect"
//	"strings"
	"syscall"
//    "unsafe"
//	"efd/efdlib"
    "golang.org/x/sys/unix"
)

const SockAddr = "/tmp/echo.sock"

func main() {
//   	var ba []byte
//    var offset  int64
//    var pid_buf [4]byte


   	conn, err := net.Dial("unix", SockAddr)
    if err != nil {
        fmt.Println("dial error:", err)
		os.Exit(1)
    }
    defer conn.Close()

	fmt.Println("client connected: ", conn.RemoteAddr().Network())

// get fd
//	fd := GetFdFromConn(conn)
//	fmt.Println(" fd: ", fd)

// alternative
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


	msg:= []byte("hello unix server")

	count, err := conn.Write(msg)
    if err != nil {
        fmt.Println("error send:", err)
		os.Exit(1)
    }

	fmt.Println("msg sent: ", string(msg), "| ", count)
// sendmsg:
// func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (err error)

	rec_msg := make([]byte,128)
	cnum :=1 // number of expected fds
	cbuf := make([]byte, syscall.CmsgSpace(cnum*4))

// func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error)
// oobn out of band control msg

	msg_num, cb_num, _, _, err := syscall.Recvmsg(int(fda), rec_msg, cbuf, 0)
	if err != nil {
		fmt.Println("error recvmsg: ", err)
		os.Exit(1)
	}

	fmt.Println("receive msg: ", string(rec_msg), " | ", msg_num)

	// parse control msgs
	var cmsgs []syscall.SocketControlMessage
	cmsgs, err = syscall.ParseSocketControlMessage(cbuf)
	if err != nil {
		fmt.Println("error parsing cntl msg: ", err)
		os.Exit(1)
	}

	fmt.Println("number of fds: ", len(cmsgs), " cb_num ", cb_num)

	shfd_ar, err := syscall.ParseUnixRights(&cmsgs[0])
	if err != nil {
		fmt.Println("error parsing cntl msg: ", err)
		os.Exit(1)
	}
	fmt.Println("efd_ar array size: ", len(shfd_ar))
	efd := shfd_ar[0]

	fmt.Println("efd: ", efd)
//    val, err = efd.ReadEvents()
	n, err := unix.Read(efd, rec_msg[:8])
    if err !=nil {
        fmt.Println("error efd read: ", err)
        return
    }
    fmt.Println("read efd: ", n, "| ", rec_msg[:8])



/*
	var filenames []string
	res := make([]*os.File, 0, len(cmsgs))

	for i := 0; i < len(cmsgs) && err == nil; i++ {
		var fds []int
		fds, err = syscall.ParseUnixRights(&cmsgs[i])

		for fi, fd := range fds {
			var filename string
			if fi < len(filenames) {
				filename = filenames[fi]
				fmt.Println("ex file: ", filename)

			}

			res = append(res, os.NewFile(uintptr(fd), filename))
			fmt.Println("file: ", filename)

		}
	}

	for i:= 0; i<len(filenames); i++ {
		fmt.Println("files: ",i, "| ", filenames[i])
	}

/*
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
*/

	msg= []byte("unix server efd received!")

	count, err = conn.Write(msg)
    if err != nil {
        fmt.Println("error send:", err)
		os.Exit(1)
    }

	fmt.Println("msg sent: ", string(msg), "| ", count)


}
