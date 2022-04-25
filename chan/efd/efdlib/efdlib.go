package evfd

/*
 * eventfd wrapper for go
 * Provides a ReadWriteCloser interface for handling eventfd()'s
 * Eventfd provides a simple filedescriptor with very low overhead.
 * It stores a bitfield of 64 bits which are added when written to
 * the fd.
 *
 * For more information on eventfd() see `man eventfd`.
 */

/*
#include <sys/eventfd.h>
*/
import "C"

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/unix"
)

type EventFD struct {
	fd    int
	valid bool
	buf [8]byte
}

//constant (
//   EFD_SEMAPHORE = 1
//    EFD_CLOEXEC = 02000000
//    EFD_NONBLOCK = 04000

const (
	EFD_CLOEXEC = unix.EFD_CLOEXEC
	EFD_SEMAPHORE = unix.EFD_SEMAPHORE
	EFD_NONBLOCK = unix.EFD_NONBLOCK
)


/* Create a new EventFD. */
func New(flag int) (*EventFD, error) {
	if flag == 0 {
		flag = EFD_CLOEXEC
	}
	Cflag := C.int(flag)
	fd, err := C.eventfd(0, Cflag)
	if err != nil {
		return nil, err
	}

	e := &EventFD{
		fd:    int(fd),
		valid: true,
	}
	return e, nil
}

/* Read events from Eventfd. p should be at max 8 bytes.
 * Returns the number of read bytes or 0 and error is set.
 */
func (e *EventFD) Read(p []byte) (int, error) {
	n, err := unix.Read(e.fd, p[:])
	if err != nil {
		return 0, err
	}
	return n, nil
}

/* Read events into a uint64 and return it. Returns 0 and error
 * if an error occured
 */
func (e *EventFD) ReadEvents() (uint64, error) {
	buf := make([]byte, 8)
	n, err := unix.Read(e.fd, buf[:])
	if err != nil {
		return 0, err
	}

	val, n := binary.Uvarint(buf)
	if n <= 0 {
		return 0, fmt.Errorf("Invalid Read")
	}

	return val, nil
}

/* Write bytes to eventfd. Will be added to the current
 * value of the internal uint64 of the eventfd().
 */
func (e *EventFD) Write(p []byte) (int, error) {
	n, err := unix.Write(e.fd, p[:])
	if err != nil {
		return 0, err
	}
	return n, nil
}

/* Write a uint64 to eventfd. Value will be added to current value
 * of the eventfd
 */
func (e *EventFD) WriteEvents(val uint64) error {
	buf := make([]byte, 8)
	n := binary.PutUvarint(buf, val)
	if n != 8 {
		return fmt.Errorf("Invalid Argument")
	}

	n, err := unix.Write(e.fd, buf[:])
	if err != nil {
		return err
	}

	if n != 8 {
		return fmt.Errorf("Could not write to eventfd")
	}

	return nil
}

/* Returns the filedescriptor which is internally used */
func (e *EventFD) Fd() int {
	return e.fd
}

/* Close the eventfd */
func (e *EventFD) Close() error {
	if e.valid == false {
		return nil
	}
	e.valid = false
	return unix.Close(e.fd)
}
