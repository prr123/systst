package main

import (
//    "net/http"
    "fmt"
	"strings"
	"os"
	"io"
	"bufio"
	"bytes"
//	"errors"
)

const ch_sz=4096

type htmlfill struct {
	ch_sz int
	infilnam string
	infil *os.File
	outfilnam string
	outfil *os.File
	buf [ch_sz]byte
}

type elchar struct {
	id string
	class string
	start int
	end int
}

type htmlpage struct {
	header elchar
	body elchar
	div []elchar
	script []elchar
	endscript []elchar
	style elchar
	endstyle elchar
}

func (html *htmlfill) open_fil() (nhtml *htmlfill, err error) {

	infilnam := html.infilnam
	fmt.Println("mod args file: ", infilnam)
// check for empty file
	infilnam = html.infilnam
	if len(infilnam) < 2 {
		return nil, fmt.Errorf("error -- filename too short! %v",err)
	}
	
// check whether filname contains .html et
	infilsplit := strings.Split(infilnam, ".")
	insplitlen := len(infilsplit)
//	fmt.Println(infilnam, "split: ",insplitlen)
	switch insplitlen {
	case 1:
		infilnam += ".html"
	case 2:
		if !strings.Contains(infilnam,".html") {
			return nil, fmt.Errorf("error -- filename has wrong extension! %v",infilsplit[1])
		}
	default:
		return nil, fmt.Errorf("error -- filename has too many dots! %v",infilnam)
	}
// check whether file + directory
	ifilstr := strings.Split(infilnam, "/")
	nisplits := len(ifilstr)
	tdir := ""
	tfilnam := ifilstr[nisplits-1]
	for i:=0; i<nisplits-1; i++ {
		tdir +=ifilstr[i]
	}

    if len(tdir) > 0 {
        info, err := os.Stat(tdir)
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("error init -- directory %v does not exist! %v", tdir, err)
        }
        if !info.IsDir() {
            return nil, fmt.Errorf("error init -- argument in cmd %v is a file not a directory! %v", err)
        }
    }

	inf, err := os.Stat(infilnam)
	if err != nil {
		fmt.Println("file: ", infilnam, "does not exist!")
		return nil, fmt.Errorf("error mod -- file %v does not exist! %v", infilnam, err)
	}
	if inf.Size() > int64(ch_sz) {
		return nil, fmt.Errorf("error mod -- file %v is too big! size: %v", infilnam, inf.Size)
	}

	infil, err := os.Open(infilnam)
    if err != nil {
        return nil, fmt.Errorf("error mod-- could not open file %v! %v", infilnam, err)
    }
	html.infil = infil
	fmt.Println("success opening file: ", infilnam)
//	fil.WriteString("</html>\n")

// create output file
	tfils := strings.Split(tfilnam, ".")
	outfilnam := tfils[0] + "new.html"

	if len(tdir) > 0 {
		outfilnam = tdir + "/" + outfilnam
	}
    _, err = os.Stat(outfilnam)

    if err == nil {
        fmt.Println("file: ", outfilnam, " already exists!")
//      return fmt.Errorf("error init -- file %v already exists! %v", filnam, err)
        fmt.Println("deleting existing file!")
        err = os.Remove(outfilnam)
        if err != nil {
                return nil, fmt.Errorf("error init -- could not delete file %v! %v", outfilnam, err)
        }
    }
    fmt.Println("creating new out file: ", outfilnam)

	outfil, err := os.Create(outfilnam)
    if err != nil {
        return nil, fmt.Errorf("error mod-- could not open file %v! %v", outfilnam, err)
    }
	html.outfil = outfil
	html.outfilnam = outfilnam
	fmt.Println("success opening file: ", outfilnam)

	nhtml = html
	return nhtml, nil
}

func (html *htmlfill) close_fil () (err error) {
	if html.infil != nil {
		err = html.infil.Close()
		if err != nil {
        	return fmt.Errorf("error mod-- could not close input file %v! %v", html.infilnam, err)
		}
		fmt.Println("closed input file!")
	}
	if html.outfil != nil {
		err = html.outfil.Close()
		if err != nil {
        	return fmt.Errorf("error mod-- could not close output file %v! %v", html.outfilnam, err)
		}
		fmt.Println("closed output file!")
	}
	return nil
}

func (html *htmlfill) read_fil () (nhtml *htmlfill, err error) {
	infil := html.infil
	rd := bufio.NewReader(infil)
//	bufptr := buf[:]
	inbufptr := html.buf[:]
	for i:=0; i<5; i++ {
		n, err := rd.Read(inbufptr)
		fmt.Println("reading chunck: ", i, " chars: ", n)
        if err != nil && err != io.EOF {
        	return nil, fmt.Errorf("error mod-- error reading input file! %v", err)
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

	nhtml = html
	return nhtml, nil
}

func (html *htmlfill) parse()(err error) {
	inbufptr := html.buf[:]
	ln := bytes.Index(inbufptr, []byte{0})
	outbufptr  := inbufptr[:ln]

	ist := 0
	for i:=0; i<100; i++ {
//		fmt.Println("ist: ", ist)
		idx := bytes.Index(outbufptr[ist:], []byte("<div"))
//		fmt.Println("idx: ", idx)
		if idx < 0 {
			break;
		}
		idx = ist+idx
		clidx := bytes.Index(outbufptr[idx:], []byte(">"))
//		fmt.Println("str index: ", idx, ":", clidx)
		fmt.Println("str: ", string(outbufptr[idx:(idx+clidx+1)]), " | ", idx, ":", clidx)
		ist = idx + clidx + 2
	}

	return nil
}

func (html *htmlfill) write_fil()(err error) {
	inbufptr := html.buf[:]
	ln := bytes.Index(inbufptr, []byte{0})
	outbufptr  := inbufptr[:ln]
	outfil := html.outfil
// need to implement chunks
	n, err := outfil.Write(outbufptr)
	if err != nil {
		return fmt.Errorf("error write -- error writing output file! %v", err)
	}
	fmt.Println("written: ", n)
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

	html := new(htmlfill)
	html.ch_sz = ch_sz
	html.infilnam = os.Args[1]
	html, err := html.open_fil()
//	bufptr, err := read_html_file(filnam, opt, bufptr)
	if err != nil {
		fmt.Println("error -- opening html file: ", err)
		os.Exit(1)
	}
	html, err = html.read_fil()
	if err != nil {
		fmt.Println("error -- opening html file: ", err)
		os.Exit(1)
	}
	bufptr := html.buf[:]
	fmt.Println("buffer: \n", string(bufptr))

//
	err = html.parse()
	err = html.write_fil()
	err = html.close_fil()
//	bufptr, err := read_html_file(filnam, opt, bufptr)
	if err != nil {
		fmt.Println("error -- closing html files: ", err)
		os.Exit(1)
	}


	fmt.Println("success -- modified html file! ")
}
