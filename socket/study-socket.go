package main

// This is a super-quick example of how to set the socket options to allow port re-use for a single address/port on a host machine.
// This is most commonly used with things like hot reloads of configuration.

import (
	"context"
	"log"
	"net"
	"syscall"
)

func main() {
	const network = "tcp"
	const address = "127.0.0.1:8080"

	config := &net.ListenConfig{Control: reusePort}

	listener1, _ := config.Listen(context.Background(), network, address) // bind to the address:port
	go listen(1, listener1)

	// as soon as listener2 is bound below, all traffic will begin to flow it instead of listener1. the only thing left to do
	// with listener1 is to shut down any active connections and close the listener.
	listener2, _ := config.Listen(context.Background(), network, address) // also bind to the address:port
	listen(2, listener2)
}

func listen(id int, listener net.Listener) {
	socket, _ := listener.Accept()
	log.Println("Accepted socket from listener", id)
	socket.Close()
	listener.Close()
}

func reusePort(network, address string, conn syscall.RawConn) error {
	return conn.Control(func(descriptor uintptr) {
		syscall.SetsockoptInt(int(descriptor), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
	})
}