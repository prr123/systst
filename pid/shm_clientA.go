package main

import (
	"fmt"
	"os"
//	"io"
	"net"
	"reflect"
	"strings"
	"syscall"
//    "unsafe"
//    "golang.org/x/sys/unix"
	shm "pid/shmlib"
)

const SockAddr = "/tmp/shm.sock"

func examiner(t reflect.Type, depth int) {
	fmt.Println(strings.Repeat("\t", depth), "Type is", t.Name(), "and kind is", t.Kind())
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		fmt.Println(strings.Repeat("\t", depth+1), "Contained type:")
		examiner(t.Elem(), depth+1)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Println(strings.Repeat("\t", depth+1), "Field", i+1, "name is", f.Name, "type is", f.Type.Name(), "and kind is", f.Type.Kind())
			if f.Tag != "" {
				fmt.Println(strings.Repeat("\t", depth+2), "Tag is", f.Tag)
				fmt.Println(strings.Repeat("\t", depth+2), "tag1 is", f.Tag.Get("tag1"), "tag2 is", f.Tag.Get("tag2"))
			}
		}
	}
}

func main() {
   	var ba []byte
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
// alternative
/*
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
*/

	fda, err := shm.ConnGetFd(conn)
	if err != nil {
        fmt.Println("error cannot obtain unix fd:", err)
		os.Exit(1)
    }
	fmt.Println(" fd alt: ", fda)

	msg:= []byte("hello unix server")

	count, err := conn.Write(msg)
    if err != nil {
        fmt.Println("error send:", err)
		os.Exit(1)
    }

	fmt.Println("msg sent: ", string(msg), "| ", count)


	rec_msg := make([]byte,128)
	cnum :=1 // number of expected fds
	cbuf := make([]byte, syscall.CmsgSpace(cnum*4))

// func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error)
// oobn out of band control msg

	msg_num, shfd, err := shm.Shm_recmsg(rec_msg, cbuf, fda)

	fmt.Println("receive msg: ", string(rec_msg), " | ", msg_num)
//	fmt.Println("number of fds: ", len(cmsgs), " cb_num ", cb_num)

/*
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
	fmt.Println("shfd array size: ", len(shfd_ar))
	shfd := shfd_ar[0]
*/


	shmem := new(shm.Shm)

	shmem.Shfd = shfd
	shmem.Size = 256
	prot := 0
	err = shmem.Open_shm(prot)
	if err != nil {
		fmt.Println("Open_shm error: ", err)
		os.Exit(1)
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
*/

	ba = shmem.Ba
	shmem.Print_shm()
    fmt.Println("memory map established!")
    fmt.Println("ba proc_b: ", string(ba))



}
