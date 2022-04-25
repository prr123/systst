package main

import (
//    "net/http"
    "fmt"
	"strings"
	"os"
	"io"
	"bufio"
//	"errors"
)

const ch_sz=4096

func mod_html_file(filnam string, opt bool) (err error) {
	var buf [ch_sz]byte
	var tdir string
	fmt.Println("mod args file: ", filnam, "del: ", opt)
// check for empty file
	if len(filnam) < 2 {
		return fmt.Errorf("error -- filename too short! %v",err)
	}
// check whether filname contains .html et
	ext_idx := strings.LastIndex(filnam, ".")
	if ext_idx < 0 {
		filnam += ".html"
	} else {
		html_idx := strings.LastIndex(filnam, ".html")
		if html_idx < 0 {
			return fmt.Errorf("error -- filename has wrong extension! %v",err)
		}
	}
// check whether file + directory
	file_idx  := strings.LastIndex(filnam, "/")
	if file_idx > 0 {
		tdir = filnam[:file_idx]
//		fmt.Println("dir: ", tdir)
// check whether filnam is directory
		info, err := os.Stat(tdir)
		if os.IsNotExist(err) {
			return fmt.Errorf("error mod -- directory %v does not exist! %v", tdir, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("error mod -- argument in cmd %v is a file not a directory! %v", err)
		}
	}


//	fmt.Println("file: ", filnam, "del: ", opt)

	_, err = os.Stat(filnam)

	if err != nil {
		fmt.Println("file: ", filnam, "does not exist!")
		return fmt.Errorf("error mod -- file %v does not exist! %v", filnam, err)
	}
	fil, err := os.Open(filnam)
    defer fil.Close()
    if err != nil {
        return fmt.Errorf("error mod-- could not open file %v! %v", filnam, err)
    }

	fmt.Println("success opening file: ", filnam)
//	fil.WriteString("</html>\n")
	rd := bufio.NewReader(fil)
	bufptr := buf[:]
	for i:=0; i<5; i++ {
		n, err := rd.Read(bufptr)
		fmt.Println("reading chunck: ", i, " chars: ", n)
        if err != nil && err != io.EOF {
        	return fmt.Errorf("error mod-- error reading file %v! %v", filnam, err)
        }
		if err == io.EOF {
			break
		}
		if n < ch_sz {
			break
		}
        if n == 0 {
            break
        }
	}

	fmt.Println("buf: \n", string(bufptr))
	return nil
}

func main() {

	arg_num := len(os.Args)
	opt := false
//	fmt.Println("args: ", arg_num)

	if arg_num < 2 {
		fmt.Println("error -- insufficient arguments!")
		fmt.Println("usage: ./huga file opt")
		os.Exit(1)
	}

	if arg_num > 3 {
		fmt.Println("error -- too many arguments!")
		fmt.Println("usage: ./huga file opt")
		os.Exit(1)
	}

	if arg_num == 3 {
		if os.Args[2] == "-o" {
			opt = true
		}
		if !opt {
		  fmt.Println("error -- invalid option argument ",os.Args[2], "! ")
		  fmt.Println("usage: ./huga file opt [-o]")
		  os.Exit(1)
		}
	}

	filnam := os.Args[1]
//	fmt.Println("from command line -- file name: ", filnam, " del: ", opt)

	err := mod_html_file(filnam, opt)
	if err != nil {
		fmt.Println("error -- modifying html file: ", err)
		os.Exit(1)
	}
	fmt.Println("success -- modified html file! ", filnam)
}
