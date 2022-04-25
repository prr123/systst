package main
import (
	"fmt"
	"os"
//	"io"
	"net"
	"reflect"
	"strings"
	"syscall"
	"log"
//    "unsafe"
//    "golang.org/x/sys/unix"
	shm "pid/shmlib"
	"strconv"
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

	msg:= []byte("hello unix server")

	count, err := conn.Write(msg)
    if err != nil {
        Logger.Fatalln("error write to server:", err)
    }

	Logger.Println("msg sent: ", string(msg), "| ", count)


	rec_msg := make([]byte,128)
	cnum :=1 // number of expected fds
	cbuf := make([]byte, syscall.CmsgSpace(cnum*4))

// func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error)
// oobn out of band control msg

	msg_num, shfd, err := shm.Shm_recmsg(rec_msg, cbuf, fda)

	Logger.Println("receive msg: ", string(rec_msg), " | ", msg_num)
//	fmt.Println("number of fds: ", len(cmsgs), " cb_num ", cb_num)

/*
	msg_num, cb_num, _, _, err := syscall.Recvmsg(int(fda), rec_msg, cbuf, 0)
	if err != nil {
		fmt.Println("error recvmsg: ", err)
		os.Exit(1)
	}

	fmt.Println("receive msg: ", string(rec_msg), " | ", msg_num)

	// parse control msgs
	var cmsgs []syscall.SocketControlMessage
	cmsgs, err = syscall.ParseSocketControlMessage(cbuf)
	if err != nil {
		fmt.Println("error parsing cntl msg: ", err)
		os.Exit(1)
	}

	fmt.Println("number of fds: ", len(cmsgs), " cb_num ", cb_num)

	shfd_ar, err := syscall.ParseUnixRights(&cmsgs[0])
	if err != nil {
		fmt.Println("error parsing cntl msg: ", err)
		os.Exit(1)
	}
	fmt.Println("shfd array size: ", len(shfd_ar))
	shfd := shfd_ar[0]
*/


	shmem := new(shm.Shm)

	shmem.Shfd = shfd
	shmem.Size = 256
	prot := 0
	err = shmem.Open_shm(prot)
	if err != nil {
		Logger.Fatalln("Open_shm error: ", err)
	}
/*
    offset = 0
    shlen := 256
    prot := unix.PROT_READ
//| unix.PROT_WRITE
    shflags:= unix.MAP_SHARED

    ba,err  = unix.Mmap(shfd, offset, shlen, prot, shflags)
    if err != nil {
        fmt.Println("Mmap error!", err)
        os.Exit(1)
    }
*/

	ba = shmem.Ba
	shmem.Print_shm("log")
    Logger.Println("memory map established!")
    Logger.Println("ba proc_b: ", string(ba))

}



func main() {

	if len(os.Args) < 2 {
		fmt.Println("Need to specify number of clients!")
		fmt.Println("Usage: shm_start_clients numclients")
		os.Exit(1)
	}

	num_clients, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("error converting argument! ", err)
		fmt.Println("Usage: shm_start_clients numclients")
		os.Exit(1)
	}

    file, err := os.OpenFile("chm_clients.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

    log.SetOutput(file)
    log.Println("Starting shm_clients")

	ab := byte(65)

	for i:=0; i< num_clients; i++ {
		ab += 1
		character := string(ab)
//		fmt.Println(i, " : ", character)
		cmd := exec.Command("shm_client", character)
    	err := cmd.Run()
    	if err != nil {
			fmt.Println(err)
			log.Fatal(err)
    	}

	} // end for loop

}
