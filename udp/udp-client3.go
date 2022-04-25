package main
import (
    "fmt"
    "net"
    "bufio"
)

func main() {
    p :=  make([]byte, 2048)
	var s string
    conn, err := net.Dial("udp", "217.69.0.247:1234")
    if err != nil {         fmt.Println("Some error %v", err)
        return
    }
	fmt.Println("connection open!")

	for i:=0; i<2; i++ {
// send message to connection
		s = fmt.Sprint("Call: ",i,": How are you doing?")
    	fmt.Fprintf(conn, s)
// read response
	    _, err = bufio.NewReader(conn).Read(p)
    	if err == nil {
			s = string(p)
        	fmt.Println("Server resp ",i,": ", s)
    	} else {
        	fmt.Println("Some error %v\n", err)
			break
    	}

    }
	fmt.Println("closing connection!")
    conn.Close()
}
