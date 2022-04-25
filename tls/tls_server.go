package main

import (
    "crypto/rand"
    "crypto/tls"
    "fmt"
    "net"
    "crypto/x509"
	"io"
	"os"
//	"time"
	"log"
)

func main() {
	logfilnam := "tls_server.log"
	_, err := os.Stat(logfilnam)
	if err == nil {
		err1 := os.Remove(logfilnam)
		if err1 != nil {
			fmt.Println("error removing logfile: ", logfilnam, " err: ", err)
			os.Exit(1)
		}
	}
	logfil, err := os.OpenFile(logfilnam, os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        fmt.Println("error creating log file: ", logfilnam, " error: ", err)
		os.Exit(1)
    }

    log.SetOutput(logfil)
	log.SetFlags(log.Lmicroseconds)

    log.Println("TLS_Server Log")


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

//		t :=time.Now()
        log.Println("server: conn accepted from: ", conn.RemoteAddr())
        tlscon, ok := conn.(*tls.Conn)
        if ok {
            log.Println("TLS ok", conn.RemoteAddr())
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
	rem := conn.RemoteAddr()

	for {
    	log.Println("server: conn: ", rem, "waiting")
        n, err := conn.Read(buf)
        if err != nil {
			if err == io.EOF {
//				t := time.Now()
        		log.Println("server ", rem, " conn closing")
				break
           	} else {
                log.Println("server ", rem, " read err: ", err)
            }
			return
        }
        log.Printf("server read: %s %q\n", conn.RemoteAddr(), string(buf[:n]))
/*
		buf = []byte("hello client")
        n, err = conn.Write(buf)
        log.Printf("server ", rem, " wrote %d bytes", buf)

        if err != nil {
            log.Printf("server ", rem, " write error: %v", err)
            break
        }
*/
    }
}
