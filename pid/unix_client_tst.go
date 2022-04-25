package main

import (
	"fmt"
	"os"
//	"io"
	"net"
	"reflect"
	"strings"
)

const SockAddr = "/tmp/echo.sock"

func examiner(t reflect.Type, depth int) {
	fmt.Println(strings.Repeat("\t", depth), "Type is", t.Name(), "and kind is", t.Kind())
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		fmt.Println(strings.Repeat("\t", depth+1), "Contained type:")
		examiner(t.Elem(), depth+1)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Println(strings.Repeat("\t", depth+1), "Field", i+1, "name is", f.Name, "type is", f.Type.Name(), "and kind is", f.Type.Kind())
			if f.Tag != "" {
				fmt.Println(strings.Repeat("\t", depth+2), "Tag is", f.Tag)
				fmt.Println(strings.Repeat("\t", depth+2), "tag1 is", f.Tag.Get("tag1"), "tag2 is", f.Tag.Get("tag2"))
			}
		}
	}
}


func GetFdFromConn(con net.Conn) int {
	
	contype := reflect.TypeOf(con)
	examiner(contype, 4)

	v := reflect.ValueOf(con)
	fmt.Println(" v: ", v.Type())
	x:=reflect.Indirect(v).FieldByName("fd")
	fmt.Println(" x: ", x.Type())

//	netFD := reflect.Indirect(reflect.Indirect(v).FieldByName("fd"))

	netFD := reflect.Indirect(x)
	examiner(netFD, 3)
	fmt.Println(" netFD ", netFD.Type())

	xfd := reflect.Indirect(netFD).FieldByName("sysfd")
	fmt.Println(" netFD2 ", xfd.Type())

	fd := int(netFD.FieldByName("sysfd").Int())
	return fd
}


func main() {
/*
    if err := os.RemoveAll(SockAddr); err != nil {
		fmt.Println("error socket removal!", err)
		os.Exit(1)
    }
*/
   	conn, err := net.Dial("unix", SockAddr)
    if err != nil {
        fmt.Println("dial error:", err)
		os.Exit(1)
    }
    defer conn.Close()

	fmt.Println("client connected: ", conn.RemoteAddr().Network())

// get fd
	fd := GetFdFromConn(conn)
	fmt.Println(" fd: ", fd)

// alternative
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


	msg:= []byte("hello unix server")

	count, err := conn.Write(msg)
    if err != nil {
        fmt.Println("error send:", err)
		os.Exit(1)
    }

	fmt.Println("msg sent: ", msg, "| ", count)
// sendmsg:
// func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (err error)


// syscall.Sendmsg(fd,
}
