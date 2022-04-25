package main

import (
	"fmt"
	"os"
//	"io"
	"net"
	"syscall"
)

const SockAddr = "/tmp/echo.sock"


func main() {
/*
    if err := os.RemoveAll(SockAddr); err != nil {
		fmt.Println("error socket removal!", err)
		os.Exit(1)
    }
*/
   	conn, err := net.Dial("unix", SockAddr)
    if err != nil {
        fmt.Println("dial error:", err)
		os.Exit(1)
    }
    defer conn.Close()

	fmt.Println("client connected: ", conn.RemoteAddr().Network())

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

	fmt.Println("msg sent: ", msg, "| ", count)

// rcvmsg:
	num_fds := 1
	msg2 := make([]byte, syscall.CmsgSpace(num_fds * 4))
	_, _, _, _, err = syscall.Recvmsg(fda, nil, msg2,0)
    if err != nil {
        fmt.Println("error msg rec", err)
		os.Exit(1)
    }

	fmt.Println("message received!")

	var msgs []syscall.SocketControlMessage
	msgs, err = syscall.ParseSocketControlMessage(msg2)
    if err != nil {
        fmt.Println("error msg rec", err)
		os.Exit(1)
    }

	num_msgs = len(msgs)

	fmt.Println("number of fds: ", num_msgs)

	fds, err := syscall.ParseUnixRights(&msgs[0])
    if err != nil {
        fmt.Println("error rights", err)
		os.Exit(1)
    }


}
