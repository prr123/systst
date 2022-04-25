package main

import (
//	"fmt"
	"os"
	"os/signal"
//	"io"
	"net"
//	"reflect"
//	"strings"
	"syscall"
	"log"
//    "unsafe"
    "golang.org/x/sys/unix"
	shm "pid/shmlib"
)

const SockAddr = "/tmp/shm.sock"

func main() {

   	var ba []byte
//    var offset  int64
//    var pid_buf [4]byte

	if len(os.Args) < 2 {
		log.Fatalln("error insufficient arguments")
	}
	clientnam := "client" + os.Args[1]

	logfilnam := "client" + os.Args[1] + ".log"

	file, err := os.OpenFile(logfilnam, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

	Logger := log.New(file, clientnam + ": ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

//    log.SetOutput(file)
    Logger.Print("Logging to file: ", logfilnam)


   	conn, err := net.Dial("unix", SockAddr)
    if err != nil {
		Logger.Fatalln("error dialling server: ", err)
    }
    defer conn.Close()

	Logger.Print("client connected: ", conn.RemoteAddr().Network())

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
        Logger.Fatalln("error cannot obtain unix fd:", err)
    }
//	fmt.Println(" fd alt: ", fda)

	pid := uint32(os.Getpid())
	ba = shm.CvtInt32toB(pid)
    Logger.Println("pid: ", pid, ba)

	sid, err := unix.Getsid(int(pid))
	if err != nil {
        Logger.Fatalln("error cannot obtain session id", err)
    }
	Logger.Println("sid: ", sid)

	pgid, err := unix.Getpgid(int(pid))
	if err != nil {
        Logger.Fatalln("error cannot obtain process group id", err)
    }
	Logger.Println("pgid: ", pgid)


	count, err := conn.Write(ba)
    if err != nil {
        Logger.Fatalln("error write to server:", err)
    }

	Logger.Println("msg sent: ", ba, "| ", count)


	rec_msg := make([]byte,128)
	cnum :=1 // number of expected fds
	cbuf := make([]byte, syscall.CmsgSpace(cnum*4))

// func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error)
// oobn out of band control msg

	msg_num, shfd, err := shm.Shm_recmsg(rec_msg, cbuf, fda)

	Logger.Println("pid: ", rec_msg[0:4], " | ", msg_num)
	server_pid := shm.CvtBtoInt32(rec_msg)

	Logger.Println("rec msg: ", server_pid, " | ", msg_num)
//	Logger.Println("receive msg: ", string(rec_msg), " | ", msg_num)

	shmem := new(shm.Shm)

	shmem.Shfd = shfd
	shmem.Size = 256
	prot := 0
	err = shmem.Open_shm(prot)
	if err != nil {
// *** also send error to server
		Logger.Fatalln("Open_shm error: ", err)
	}

	ba = shmem.Ba
//	shmem.Print_shm(0,0)
    Logger.Println("memory map established!")
    Logger.Println("ba shared mem: ", string(ba))



    c := make(chan os.Signal, 1)
    signal.Notify(c)
	err = unix.Setpgid(int(pid), int(server_pid))
    if err != nil {
        Logger.Fatalln("error setting group id:", err)
    }

// send ok
	cmsg := []byte("ok")
	count, err = conn.Write(cmsg)
    if err != nil {
        Logger.Fatalln("error writing ok to server:", err)
    }

// wait for signal from server

    Logger.Println("waiting for signal!")

    // Block until any signal is received.
    s := <-c
    Logger.Println("Got signal:", s)

}
