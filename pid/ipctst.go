//
package main

import (
	"fmt"
//	ipc "github.com/aleccunningham/posixipc"
	ipc "pid/posixipc"
)

func main() {
	fmt.Print("Hello")

	ipc := ipc.NewOperator()

	messages := make(<-chan struct{})
	for i := range hot {
		messages <- i
	}

	ipc.Start(messages)
}
