package main

import (
	"fmt"
	"os"
//	"io"
	"net"
	"strconv"
//	"syscall"
//    "unsafe"
//  	"golang.org/x/sys/unix"
	shm  "pid/shmlib"
)

const SockAddr = "/tmp/shm.sock"

func main() {
	var ba []byte
//	var pid_ba [4]byte
//   var offset  int64

	if len(os.Args) < 2 {
		fmt.Println("no client number specified!")
		fmt.Println("Usage is: shm_server2 clients")
		os.Exit(1)
	}

	max_clients, err :=	strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("conversion error: ", err)
		os.Exit(1)
	}

	fmt.Println("max number of clients: ", max_clients)

// get process id
	pid := uint32(os.Getpid())
//	ba = pid_ba[0:3]
	ba = shm.CvtInt32toB(pid)
	fmt.Println("pid: ", pid, ba)

// create shm structure
	shmem := new(shm.Shm)
	shmem.Name = "shprr"
	shmem.Size = 256

	err = shmem.Init_shm()
	if err != nil {
		fmt.Println("Shared memory init error: ", err)
		os.Exit(1)
	}

// check protection
	prot := 0
	err = shmem.Open_shm(prot)
	if err != nil {
		fmt.Println("Shared memory open error: ", err)
		os.Exit(1)
	}

	shmem.Print_shm("hello")

    copy(shmem.Ba, "hello shared memory opened!")

//    fmt.Println(" ba: ", ba)
	fmt.Println("ba: ", string(shmem.Ba))

// preemptively clear all existing sockets
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

    for i:=0; i<max_clients; i++ {
        // Accept new connections
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("accept error:", err)
			os.Exit(1)
        }
		fmt.Println("client connected: ", conn.RemoteAddr().Network())

// a go function should start here instead of inline code

		rec_msg := make([]byte,128)
		count, err := conn.Read(rec_msg)
        if err != nil {
            fmt.Println("read error:", err)
			os.Exit(1)
        }

// replace with pid of client process

		fmt.Println("pid: ", rec_msg[0:3], " | ", count)
		client_pid := shm.CvtBtoInt32(rec_msg)

		fmt.Println("client pid: ", client_pid, " | ", count)


		cfd , err := shm.ConnGetFd(conn)
		if err != nil {
			fmt.Println("error obtaining file descriptor from conn!", err)
			os.Exit(1)
		}
//		fmt.Println(" cfd : ", cfd)

// Send the client the pid of the server process

//		nmsg := []byte("client: msg received!")
		shfd := shmem.Shfd

// may need to expand to send pointers to multiple shared file areas
// ba is pid of this process
		err = shm.Shm_sendmsg(ba, int(cfd), shfd)
        if err != nil {
            fmt.Println("sendmsg error:", err)
			os.Exit(1)
        }
// confirm that client has received pid and was able to open shared memory
		count, err = conn.Read(rec_msg)
        if err != nil {
            fmt.Println("read error:", err)
			os.Exit(1)
        }

		fmt.Println("ok msg: ", rec_msg)

		fmt.Println("closing!")
		conn.Close()
	} // for

	err = shmem.Shm_umap()
	if err != nil {
		fmt.Println("shared memory unmap error:", err)
		os.Exit(1)
	}

}
