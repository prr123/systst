package main

import (
	"fmt"
	"os"
//	"io"
	"net"
//	"strconv"
//	"syscall"
//    "unsafe"
//  	"golang.org/x/sys/unix"
	shm  "pid/shmlib"
)

const SockAddr = "/tmp/shm.sock"

func main() {
//   var ba []byte
//   var offset  int64

	var cons [10]net.Conn

	shmem := new(shm.Shm)
	shmem.Name = "shprr"
	shmem.Size = 256

	err := shmem.Init_shm ()
	if err != nil {
		fmt.Println("Shared memory init error: ", err)
		os.Exit(1)
	}

	prot := 0
	err = shmem.Open_shm(prot)
	if err != nil {
		fmt.Println("Shared memory open error: ", err)
		os.Exit(1)
	}

	shmem.Print_shm()
    copy(shmem.Ba, "hello shared memory opened!")

//    fmt.Println(" ba: ", ba)
	fmt.Println("ba: ", string(shmem.Ba))


    if err := os.RemoveAll(SockAddr); err != nil {
		fmt.Println("error socket removal!", err)
		os.Exit(1)
    }

    l, err := net.Listen("unix", SockAddr)
    if err != nil {
        fmt.Println("listen error:", err)
		os.Exit(1)
    }
/*
	ulis:= l.(*net.UnixListener)

	ufil, err :=ulis.File()
    if err != nil {
        fmt.Println("file error:", err)
		os.Exit(1)
    }

	fd := ufil.Fd()
	nam := ufil.Name()

	fmt.Println(" fd: ", strconv.Itoa(int(fd)), " Name: ", nam)
*/
    defer l.Close()

	fmt.Println("unix server listening:")

	rec_msg := make([]byte,128)

    for i:=0; i<5; i++ {
        // Accept new connections, dispatching them to echoServer
        // in a goroutine.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("accept error:", err)
			os.Exit(1)
        }
		cons[i] = conn
		fmt.Println("client connected: ", conn.RemoteAddr().Network())

		count, err := conn.Read(rec_msg)
        if err != nil {
            fmt.Println("read error:", err)
			os.Exit(1)
        }

		fmt.Println("msg : ", string(rec_msg), " | ", count)


		cfd , err := shm.ConnGetFd(conn)
		if err != nil {
			fmt.Println("error obtaining file descriptor from conn!", err)
			os.Exit(1)
		}
		fmt.Println(" cfd : ", cfd)

/*
// need to get the fd of the shred file area
    flag:= os.O_RDWR | os.O_CREATE
    shfilnam := "/dev/shm/prr"

    shfil, err := os.OpenFile(shfilnam, flag, 0667)
    if err != nil {
        fmt.Println("shared file open error: ", err)
        os.Exit(1)
     }
    fmt.Println("shared file: ", shfil)

//    shfda := int(shfil.Fd())

*/

		nmsg := []byte("client: msg received!")
		shfd := shmem.Shfd

		err = shm.Shm_sendmsg(nmsg, int(cfd), shfd)
        if err != nil {
            fmt.Println("sendmsg error:", err)
			os.Exit(1)
        }

		fmt.Println("closing!")
		conn.Close()
	} // for
}
