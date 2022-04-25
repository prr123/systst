//https://github.com/gxed/eventfd/tree/80a92cca79a8041496ccc9dd773fcb52a57ec6f9
package main

import (
	"fmt"
	"github.com/sahne/eventfd"
)

func main() {
	efd, err := eventfd.New()
	if err != nil {
		fmt.Println("Could not create EventFD: %v", err)
		return
	}
	/* TODO: register fd at kernel interface (for example cgroups memory watcher) */
	/* listen for new events */
	for {
		val, err := efd.ReadEvents()
		if err != nil {
			fmt.Printf("Error while reading from eventfd: %v", err)
			break
		}
		fmt.Printf("value: ", val)
	}
}
