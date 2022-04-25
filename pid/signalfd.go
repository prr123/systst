//https://gist.github.com/paulzhol/4f984aacacb753362c943cd5f8cf85f3
package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func testSignalFd() {
	var sigmask unix.Sigset_t
	fd, err := unix.Signalfd(-1, &sigmask, unix.SFD_CLOEXEC|unix.SFD_NONBLOCK)
	if err != nil {
		log.Fatalf("failed to create a signalfd: %v", err)
	}
	signalfd := os.NewFile(uintptr(fd), "signalfd")
	defer signalfd.Close()

	sigmask.Val[0] |= 1 << uint(unix.SIGUSR1-1)
	_, err = unix.Signalfd(fd, &sigmask, 0)
	if err != nil {
		log.Fatalf("failed to update sigmask on signalfd: %v", err)
	}

	log.Printf("Sending SIGUSR1 in 2 sec")
	t := time.AfterFunc(2 * time.Second, func() {
		unix.Kill(unix.Getpid(), unix.SIGUSR1)
	})
	defer t.Stop()

	log.Printf("Waiting for SIGUSR1 from signalfd")
	buf := make([]byte, unsafe.Sizeof(unix.SignalfdSiginfo{}))
	n, err := signalfd.Read(buf)
	if err != nil {
		log.Fatalf("failed reading siginfo from signalfd: %v", err)
	}
	info := *(*unix.SignalfdSiginfo)(unsafe.Pointer(&buf[0]))
	log.Printf("n: %d, %+v", n, info)
}

func main() {
	log.SetFlags(0)
	maskExec := flag.Bool("mask_exec", false, "mask SIGUSR1 and exec self")
	flag.Parse()

	if *maskExec {
		runtime.LockOSThread()
		const SIG_SETMASK = 2
		var sigmask unix.Sigset_t
		sigmask.Val[0] |= 1 << (uint(unix.SIGUSR1) - 1)
		unix.Syscall6(unix.SYS_RT_SIGPROCMASK, SIG_SETMASK, uintptr(unsafe.Pointer(&sigmask)), 0 /*_NSIG/8 */, 8, 0, 0)
		var execPath [256]byte
		n, err := unix.Readlink("/proc/self/exe", execPath[:])
		if err != nil {
			log.Fatalf("failed to read /proc/self/exe symlink: %v", err)
		}
		err = unix.Exec(string(execPath[:n]), []string{string(execPath[:n])}, os.Environ())
		log.Fatal(err)
	}

	testSignalFd()
}