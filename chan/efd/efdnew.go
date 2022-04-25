//https://github.com/gxed/eventfd/tree/80a92cca79a8041496ccc9dd773fcb52a57ec6f9
package main

import (
	"fmt"
//	"github.com/sahne/eventfd"
	"efd/efdlib"
	"syscall"
	"os"
	"strconv"
	"net"
	"golang.org/x/sys/unix"
)

const SockAddr = "/tmp/echo.sock"

func main() {
	var ba []byte

	var val uint64

	flag := evfd.EFD_CLOEXEC | evfd.EFD_SEMAPHORE
	efd, err := evfd.New(flag)
	if err != nil {
		fmt.Println("Could not create EventFD: %v", err)
		return
	}

	eFD := int(efd.Fd())
	fmt.Println("fd: ", eFD)

	ba = []byte{4,0,0,0,0,0,0,0}
	pid := syscall.Getpid()


	fmt.Println("pid: ", pid," ba: ", ba, len(ba))

	n, err := unix.Write(eFD, ba)
	if err != nil {
		fmt.Println("error efd syscall write: ", err)
		return
	}
	fmt.Println("write:  ", n)


	val = 0
/*
	err = efd.WriteEvents(val)
	if err !=nil {
		fmt.Println("error efd write: ", err)
		return
	}

	fmt.Println("write:  ", val)
*/
	val, err = efd.ReadEvents()
	if err !=nil {
		fmt.Println("error efd read: ", err)
		return
	}
	fmt.Println("read efd: ", val)

	efd_counter := make([]byte, 8)
	n,err = unix.Read(eFD,efd_counter)
	fmt.Println ("efd counter: ",  n, "| ", efd_counter)
	/* TODO: register fd at kernel interface (for example cgroups memory watcher) */
	/* listen for new events */

/*
	for {
		val, err := efd.ReadEvents()
		if err != nil {
			fmt.Printf("Error while reading from eventfd: %v", err)
			break
		}
		fmt.Printf("value: ", val)
	}
*/

//   var ba []byte
//   var offset  int64

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
	
	fmt.Println("rec msg: ", count, "| ", string(rec_msg))

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

	nmsg := []byte("TO client: msg received!")
//      msg_ar[0] = nmsg

	rights := syscall.UnixRights(eFD)
        err = syscall.Sendmsg(int(fda), nmsg, rights, nil, 0)
        if err != nil {
            fmt.Println("sendmsg error:", err)
            os.Exit(1)
        }

        fmt.Println("closing!")

	val, err = efd.ReadEvents()
	if err !=nil {
		fmt.Println("error efd read: ", err)
		return
	}
	fmt.Println("read efd: ", val)

        count, err = conn.Read(rec_msg)
        if err != nil {
            fmt.Println("read error:", err)
            os.Exit(1)
        }
	fmt.Println("rec msg: ", count, "| ", string(rec_msg))

	efd.Close()
	conn.Close()

}

