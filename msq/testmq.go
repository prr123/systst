package main

/*
#include <stdlib.h>
#include <fcntl.h>
#include <signal.h>
#include <mqueue.h>
mqd_t mq_open4(const char *name, int oflag, int mode, struct mq_attr *attr) {
    return mq_open(name, oflag, mode, attr);
}
*/
import "C"

import (
	"fmt"
//	"unsafe"
)

func main() {

//	var md C.struct_mqd_t

	 ma := C.struct_mq_attr {
		mq_flags: C.long(0),                // blocking read/write
    	mq_maxmsg: C.long(2),              // maximum number of messages allowed in queue
    	mq_msgsize: C.long(64),    // messages are contents of an int
    	mq_curmsgs: C.long(0),              // number of messages currently in queue
	}

//	name := "/test_queue"
//	nam_ptr := (*C.uchar)(unsafe.Pointer(&buf.Bytes()[0]))

   	name := C.CString("/test_queue")
//    defer C.free(name)


	md, err := C.mq_open4(name, C.O_RDWR | C.O_CREAT, C.int(0666), &ma)
	if err != nil {
		fmt.Println("error creating queue!", err)
	}

	fmt.Println("success!", md)
}
