package shm
// vers 1.1 2/2/2021 add errors

// todo
// 1. unix.mlock
// 2.  unix.Flock
// 3. unix.mprotect
// 4. unix.FcntlInt
// 5. unix.FcntlFlock
// func FcntlInt(fd uintptr, cmd, arg int) (int, error) {
// FcntlFlock performs a fcntl syscall for the F_GETLK, F_SETLK or F_SETLKW command.
// https://man7.org/linux/man-pages/man2/fcntl.2.html
// func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error {
// log

// 19/12/2020 added pagesize

import (
	"fmt"
	"os"
	"net"
	"syscall"
  	"golang.org/x/sys/unix"
	"unsafe"
//	"errors"
//	"github.com/pkg/errors"
)

const SockAddr = "/tmp/shm.sock"

const (
    PROT_READ = unix.PROT_READ
    PROT_WRITE = unix.PROT_WRITE
	PROT_EXEC = unix.PROT_EXEC
	PROT_GROWSDOWN = unix.PROT_GROWSDOWN
	PROT_GROWSUP = unix.PROT_GROWSUP
	PROT_NONE = unix.PROT_NONE
)

type Shm struct {
	Name string
	Page int
	Pgnum int
	Size int
	Shfd int
	Prot int
	Ba []byte
}

type Usock struct {
	Sockadr string
	Ufd int
	Ucon net.Conn
}

func CvtBtoInt32(b []byte) uint32 {
    // equivalnt of return int32(binary.LittleEndian.Uint32(b))
    return uint32(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24)
}

func CvtInt32toB(h uint32)(b []byte) {
    b =(*[4]byte)(unsafe.Pointer(&h))[:]
    return
}

func CvtBtoInt64(b []byte) uint64 {
    return uint64(uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<54)
}

func CvtInt64toB(h uint64)(b []byte) {
    b = (*[8]byte)(unsafe.Pointer(&h))[:]
    return
}

/*
func toByteArray(i int32) (arr [4]byte) {
    binary.BigEndian.PutUint32(arr[0:4], uint32(i))
    return
}
*/

func ConvInt32toBa(i int32) (arr[4]byte) {
	*(*int32)(unsafe.Pointer(&arr[0])) = i
	return
}

func ConvBatoInt32(arr [4]byte) (i_ret int32) {
	i_ret = *(*int32)(unsafe.Pointer(&arr[0]))
	return i_ret
}

func ConvInt64toBa(i int64) (arr[8]byte) {
	*(*int64)(unsafe.Pointer(&arr[0])) = i
	return
}

func ConvBatoInt64(arr [8]byte) (i_ret int64) {
	i_ret = *(*int64)(unsafe.Pointer(&arr[0]))
	return i_ret
}


func (shm *Shm) Prot_shm(prot ...int) {
	tprot := 0
	for iprot := range prot {
		tprot = tprot | iprot
	}
	shm.Prot = tprot
	return
}

func SockFd (l net.Listener) (ufil *os.File, err error) {

	ulis := l.(*net.UnixListener)
	ufil, err = ulis.File()
    if err != nil {
        return nil, fmt.Errorf("get SockFd error: %w ", err)
    }
	return ufil, nil
}

func SockFD (ufil *os.File) (fd int) {
	fd = int(ufil.Fd())
	return fd
}

func SockName(ufil *os.File) (nam string) {
	nam = ufil.Name()
	return nam
}

func ConnGetFd (conn net.Conn) (fda int, err error) {
    ucon, ok := conn.(*net.UnixConn)
    if !ok {
       	return -1, fmt.Errorf("error: could not establish unix socket connection")
    }
    file, err := ucon.File()
    if err != nil {
//        fmt.Println("ucon file error:", err)
        return -1, fmt.Errorf("error ucon file retrieval: %w ", err)
    }

    fda =  int(file.Fd())
	return fda, nil
}

func (shm *Shm) Print_shm (st, end int) {
// still todo
	fmt.Println("************** shm print ***********")
	if shm == nil {
		fmt.Println("shm struct not initialized")
		return
	}
	fmt.Println("Name:   ", shm.Name, "| len: ", len(shm.Name))
	fmt.Println("Size:   ", shm.Size)
	fmt.Println("Page:   ", shm.Page)
	fmt.Println("Num_pg: ", shm.Pgnum)
	fmt.Println("Fd:     ", shm.Shfd)
	fmt.Println("Buf[", st, ":", end,"]: ", shm.Ba[st:end], "| Len: ", len(shm.Ba))
	fmt.Println("************** shm print ***********")
	return
}

func (shm *Shm) Init_shm () (err error){

	if shm == nil {
		return fmt.Errorf("Init_shm error: need valid Shm structure")
	}

	name := shm.Name
	dsize := shm.Size

// find the number of pages
	pagesize := os.Getpagesize()
	shm.Page = pagesize
	num_pages := int(dsize/pagesize)
	if dsize > num_pages*pagesize {
		num_pages++
	}
	shm.Pgnum = num_pages
	size := num_pages*pagesize
	shm.Size = size

	if len(name) < 4 {
		return fmt.Errorf("Init_shm error: need name longer than 4 chars")
	}
	if size < 64 {
		return fmt.Errorf("Init_shm error: need size larger than 64 bytes")
	}

	shflags := unix.MFD_CLOEXEC | unix.MFD_ALLOW_SEALING
// **
	shfd, err := unix.MemfdCreate(name, shflags)
    if err != nil {
        return fmt.Errorf("shared memory fd creation error! %w", err)
    }

	shm.Shfd = shfd

    err = unix.Ftruncate(shfd, int64(size))
    if err != nil {
        return fmt.Errorf("ftruncate error! %w", err)
    }
	return nil
}

func (shm *Shm) Open_shm  (prot int) (err error) {

    offset := int64(0)
	shlen := shm.Size
	shfd := shm.Shfd
    prot = unix.PROT_READ | unix.PROT_WRITE
    shflags := unix.MAP_SHARED
    ba, err  := unix.Mmap(shfd, offset, shlen, prot, shflags)
    if err != nil {
        return fmt.Errorf("Mmap error! %w", err)
    }
	shm.Ba = ba
	return nil
}

func (shm *Shm) Lock_shm () (err error) {

	err = unix.Mlock(shm.Ba)
    if err != nil {
        return fmt.Errorf("Mlock error! %w", err)
    }
	return nil
}

func Shm_sendmsg(msg []byte, fda int, shfd int)(err error) {
// shared memory file descriptor Shm.shfd
// unix socket file descriptor: fda
	rights := syscall.UnixRights(shfd)
    err = syscall.Sendmsg(int(fda), msg, rights, nil, 0)
	if err != nil {
        return fmt.Errorf("sendmsg error: %w", err)
	}
	return nil
}

func Shm_recmsg(rec_msg, cbuf []byte, fda int)(msg_num, cfd int, err error) {

	var cmsgs []syscall.SocketControlMessage

	cnum := 1
	cbuf = make([]byte, syscall.CmsgSpace(cnum*4))
//	msg_num, cb_num, _, _, err := syscall.Recvmsg(int(fda), rec_msg, cbuf, 0)
	msg_num, _, _, _, err = syscall.Recvmsg(int(fda), rec_msg, cbuf, 0)
    if err != nil {
        return 0, -1, fmt.Errorf("error recvmsg: %w", err)
    }

    cmsgs, err = syscall.ParseSocketControlMessage(cbuf)
    if err != nil {
        return 0, -1, fmt.Errorf("error parsing cntl msg: %w", err)
    }

//    fmt.Println("number of fds: ", len(cmsgs), " cb_num ", cb_num)

    shfd_ar, err := syscall.ParseUnixRights(&cmsgs[0])
    if err != nil {
        return 0, -1, fmt.Errorf("error parsing cntl msg: %w", err)
    }
//    fmt.Println("shfd array size: ", len(shfd_ar))
    shfd := shfd_ar[0]

	return msg_num, shfd, nil
}

// code to expand for receipt of multiple file descriptors
/*
    var filenames []string
    res := make([]*os.File, 0, len(cmsgs))

    for i := 0; i < len(cmsgs) && err == nil; i++ {
        var fds []int
        fds, err = syscall.ParseUnixRights(&cmsgs[i])

        for fi, fd := range fds {
            var filename string
            if fi < len(filenames) {
                filename = filenames[fi]
                fmt.Println("ex file: ", filename)

            }

            res = append(res, os.NewFile(uintptr(fd), filename))
            fmt.Println("file: ", filename)

        }
    }

    for i:= 0; i<len(filenames); i++ {
        fmt.Println("files: ",i, "| ", filenames[i])
    }
*/



func (shm *Shm) Shm_umap () (err error) {
	ba := shm.Ba
	err = unix.Munmap(ba)
	if err != nil {
		return fmt.Errorf("unmap error %w", err)
	}
	return nil
}

