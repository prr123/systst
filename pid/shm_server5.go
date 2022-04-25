package main

// version 5
// create exit function to exit with unmmaping the shared memory
// 1. changed file -> logfile


import (
	"fmt"
	"os"
	"log"
//	"io"
	"net"
	"strconv"
//	"syscall"
//    "unsafe"
  	"golang.org/x/sys/unix"
	shm  "pid/shmlib"
)

const SockAddr = "/tmp/shm.sock"

var Tstlog *log.Logger
//	shmem := new(shm.Shm)
var shmem shm.Shm

func main() {
	var pid_ba, sid_ba []byte
//	var ba [O]byte
//   var offset  int64
   	var tstsig unix.Signal

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

	logfilnam := "shm_server5.log"


//	file, err := os.OpenFile(logfilnam, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	logfile, err := os.OpenFile(logfilnam, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    defer logfile.Close()

    Tstlog = log.New(logfile, "server: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

//    log.SetOutput(file)
    Tstlog.Print("Logging to file: ", logfilnam)
	Tstlog.Print("max number of clients: ", max_clients)

// get process id
	pid := uint32(os.Getpid())
//	ba = pid_ba[0:3]
	pid_ba = shm.CvtInt32toB(pid)
	Tstlog.Print("server pid: ", pid)

// session id
	sid, err := unix.Getsid(int(pid))
	if err != nil {
		Tstlog.Fatalln("Error getting session id: ", err)
	}
	sid_ba = shm.CvtInt32toB(uint32(sid))
	Tstlog.Print("session id: ", sid)


// create shm structure
//	shmem := new(shm.Shm)
	shmem.Name = "shprr"
	shmem.Size = 4096

	err = shmem.Init_shm()
	if err != nil {
		Tstlog.Fatalln("Shared memory init error: ", err)
	}

// check protection
	prot := 0
	err = shmem.Open_shm(prot)
	if err != nil {
		Tstlog.Fatalln("Shared memory open error: ", err)
	}

	for i:=0; i<32; i++ {
		shmem.Ba[i] = 0
	}
    copy(shmem.Ba[8:11],pid_ba)
    copy(shmem.Ba[12:15],sid_ba)


//	shmem.Print_shm(0,16)


//    fmt.Println(" ba: ", ba)
	Tstlog.Print("shmem[0:16]: ", shmem.Ba[0:16])

// preemptively clear all existing sockets
    if err := os.RemoveAll(SockAddr); err != nil {
		exit_main("error socket removal!", err)
    }

    l, err := net.Listen("unix", SockAddr)
    if err != nil {
		exit_main("socket listen error:", err)
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

	Tstlog.Print("unix server listening:")

    for i:=0; i<max_clients; i++ {
        // Accept new connections
        conn, err := l.Accept()
        if err != nil {
			exit_main("socket accept error:", err)
        }
		Tstlog.Print("client connected: ", conn.RemoteAddr().Network())

// a go function should start here instead of inline code

		rec_msg := make([]byte,128)
		count, err := conn.Read(rec_msg)
        if err != nil {
			exit_main("rec_msg read error:", err)
        }

// replace with pid of client process

//		Tstlog.Print("client pid: ", rec_msg[0:3], " | ", count)
		client_pid := shm.CvtBtoInt32(rec_msg)

		Tstlog.Print("client pid: ", client_pid, " | ", count)


		cfd , err := shm.ConnGetFd(conn)
		if err != nil {
			exit_main("error getting fd from conn!", err)
		}
//		fmt.Println(" cfd : ", cfd)

// Send the client the pid of the server process

//		nmsg := []byte("client: msg received!")
		shfd := shmem.Shfd

// may need to expand to send pointers to multiple shared file areas
// ba is pid of this process

		err = shm.Shm_sendmsg(pid_ba, int(cfd), shfd)
        if err != nil {
			exit_main("sendmsg error:", err)
        }
// confirm that client has received pid and was able to open shared memory
		count, err = conn.Read(rec_msg)
        if err != nil {
			exit_main("rec_msg 2 error:", err)
        }

		if count == 0 {
			exit_main("no client ok msg received", nil)
		} else {

			Tstlog.Print("client ok msg: ", string(rec_msg))
			rec_str := string(rec_msg[:2])
			if rec_str !=  "ok" {
				s := "error received message not ok: " + string(rec_msg[0:count])
				exit_main(s, nil)
			}
// need to create an exit routine
		}

    	tstsig = 34

    	err = unix.Kill(int(client_pid), tstsig)
    	if err != nil {
			s2 := fmt.Sprintf("error sending signal to: %d error: ", client_pid)
        	exit_main(s2, err)
    	}
		Tstlog.Print("sent signal ", tstsig, " to: ", client_pid)
		Tstlog.Print("closing unix connection!", i)
		conn.Close()

	} // for

	Tstlog.Print("Exiting Server! -- unmap")
	err = shmem.Shm_umap()
	if err != nil {
		Tstlog.Fatalln("shared memory unmap error:", err)
	}
	Tstlog.Print("Exiting Server! -- success")
	os.Exit(0)
}

func exit_main(msg string, err error) {
	Tstlog.Print(msg + " %w", err)
	Tstlog.Print("Fatal Error Exiting - unmapping shared memory")
	err = shmem.Shm_umap()
	if err != nil {
		Tstlog.Print("shared memory unmap error: %w", err)
	}
	os.Exit(1)
}

