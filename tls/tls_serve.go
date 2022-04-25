package main

import (
    "crypto/rand"
    "crypto/tls"
    "fmt"
    "net"
    "crypto/x509"
	"io"
	"os"
	"time"

)

func main() {
    cert, err := tls.LoadX509KeyPair("/home/peter/newca/server_crt.pem", "/home/peter/newca/server_key.pem")
    if err != nil {
        fmt.Println("server: loadkeys: ", err)
		os.Exit(1)
    }
    config := tls.Config{Certificates: []tls.Certificate{cert}}
    config.Rand = rand.Reader

    service := "127.0.0.1:8011"
    listener, err := tls.Listen("tcp", service, &config)
    if err != nil {
        fmt.Println("listener create error: ", err)
		os.Exit(1)
    }
    fmt.Println("server: listening")
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("server: accept error: ", err)
            break
        }
        defer conn.Close()

		t :=time.Now()
        fmt.Println("server: conn accepted at:", t.Format("15:04:05.000000"), " from ", conn.RemoteAddr())
        tlscon, ok := conn.(*tls.Conn)
        if ok {
            fmt.Println("ok=true")
            state := tlscon.ConnectionState()
            for _, v := range state.PeerCertificates {
                fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
            }
        }
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 512)
    for {
        fmt.Println("server: conn: ", conn.RemoteAddr(), "waiting")
        n, err := conn.Read(buf)
        if err != nil {
			if err == io.EOF {
				t := time.Now()
        		fmt.Println("server: conn closing at:", t.Format("15:04:05.000000"), " from ", conn.RemoteAddr())
				break
           	} else {
                fmt.Println("server: conn: read: ", err)
            }
			return
        }
        fmt.Printf("server read: %s %q\n", conn.RemoteAddr(), string(buf[:n]))
/*
        n, err = conn.Write(buf[:n])
        log.Printf("server: conn: wrote %d bytes", n)

        if err != nil {
            log.Printf("server: write: %s", err)
            break
        }
*/
    }
    fmt.Println("server: conn: closed")
}
