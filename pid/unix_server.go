package main

import (
	"fmt"
	"os"
//	"io"
	"net"
)

const SockAddr = "/tmp/echo.sock"

func main() {

    if err := os.RemoveAll(SockAddr); err != nil {
		fmt.Println("error socket removal!", err)
		os.Exit(1)
    }

    l, err := net.Listen("unix", SockAddr)
    if err != nil {
        fmt.Println("listen error:", err)
		os.Exit(1)
    }
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

		fmt.Println("closing!")
		conn.Close()
//	} // for
}
